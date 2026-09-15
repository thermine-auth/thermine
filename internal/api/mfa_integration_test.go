package api

import (
	"net/http"
	"testing"
	"time"

	"xermess/internal/config"
	"xermess/internal/totp"
)

// codeFor is the authenticator's code for a moment.
func codeFor(t *testing.T, secret string, at time.Time) string {
	t.Helper()
	code, err := totp.Code(secret, totp.Step(at))
	if err != nil {
		t.Fatal(err)
	}
	return code
}

type loginAnswer struct {
	Next  string         `json:"next"`
	Admin map[string]any `json:"admin"`
}

func (c *client) loginAnswer(email, password string) (int, loginAnswer) {
	var out loginAnswer
	status := c.do(http.MethodPost, "/auth/login", map[string]string{"username": email, "password": password}, &out)
	return status, out
}

func (c *client) state() string {
	var out struct {
		State string `json:"state"`
	}
	c.must(http.StatusOK, http.MethodGet, "/auth/session", nil, &out)
	return out.State
}

// Where two-factor sign-in is required, an administrator without it is made to
// set it up before anything else, and signs in with it from then on.
func TestLiveAdminMFARequired(t *testing.T) {
	s := newLiveServerWith(t, func(cfg *config.Config) { cfg.AdminMFARequired = true })

	root := s.client()
	root.must(http.StatusCreated, http.MethodPost, "/setup", map[string]string{
		"email": superEmail, "password": superPassword, "first_name": "Root",
	}, nil)

	// The password alone only gets as far as setting up.
	status, answer := root.loginAnswer(superEmail, superPassword)
	if status != http.StatusOK || answer.Next != "enroll" || answer.Admin != nil {
		t.Fatalf("login = %d %+v, want next=enroll and no admin", status, answer)
	}
	if got := root.state(); got != "enroll" {
		t.Errorf("state = %q, want enroll", got)
	}
	for _, path := range []string{"/me", "/users", "/admins"} {
		if status := root.do(http.MethodGet, path, nil, nil); status != http.StatusUnauthorized {
			t.Errorf("GET %s before setting up = %d, want 401", path, status)
		}
	}

	var begun struct {
		Enrolment struct {
			Secret string `json:"secret"`
			URI    string `json:"uri"`
		} `json:"enrolment"`
	}
	root.must(http.StatusOK, http.MethodPost, "/mfa/totp", nil, &begun)
	secret := begun.Enrolment.Secret
	if secret == "" || begun.Enrolment.URI == "" {
		t.Fatalf("enrolment = %+v", begun)
	}

	if status := root.do(http.MethodPost, "/mfa/totp/confirm", map[string]string{"code": "000000"}, nil); status != http.StatusBadRequest {
		t.Errorf("confirming with a wrong code = %d, want 400", status)
	}

	var confirmed struct {
		RecoveryCodes []string `json:"recovery_codes"`
	}
	now := time.Now()
	root.must(http.StatusOK, http.MethodPost, "/mfa/totp/confirm", map[string]string{"code": codeFor(t, secret, now)}, &confirmed)
	if len(confirmed.RecoveryCodes) != 10 {
		t.Fatalf("recovery codes = %v, want 10", confirmed.RecoveryCodes)
	}

	// Confirming signed the waiting session in.
	var me struct {
		Admin struct {
			MFAEnabled bool `json:"mfa_enabled"`
		} `json:"admin"`
	}
	root.must(http.StatusOK, http.MethodGet, "/me", nil, &me)
	if !me.Admin.MFAEnabled {
		t.Error("me.mfa_enabled = false after confirming")
	}

	// It cannot be turned off where it is required.
	if status := root.do(http.MethodDelete, "/mfa/totp", map[string]string{"code": confirmed.RecoveryCodes[9]}, nil); status != http.StatusConflict {
		t.Errorf("turning off a required factor = %d, want 409", status)
	}

	// Next sign-in: a code is needed, and the used step cannot be replayed.
	laptop := s.client()
	if _, answer := laptop.loginAnswer(superEmail, superPassword); answer.Next != "mfa" {
		t.Fatalf("second login next = %q, want mfa", answer.Next)
	}
	if status := laptop.do(http.MethodGet, "/me", nil, nil); status != http.StatusUnauthorized {
		t.Errorf("me waiting for a code = %d, want 401", status)
	}
	// A session that has the password but not the phone cannot enrol its own
	// authenticator.
	if status := laptop.do(http.MethodPost, "/mfa/totp", nil, nil); status != http.StatusUnauthorized {
		t.Errorf("starting a set-up while waiting for a code = %d, want 401", status)
	}

	if status := laptop.do(http.MethodPost, "/auth/mfa", map[string]string{"code": codeFor(t, secret, now)}, nil); status != http.StatusUnauthorized {
		t.Errorf("replaying the code used to confirm = %d, want 401", status)
	}

	next := codeFor(t, secret, now.Add(totp.Period))
	laptop.must(http.StatusOK, http.MethodPost, "/auth/mfa", map[string]string{"code": next}, nil)
	laptop.must(http.StatusOK, http.MethodGet, "/me", nil, nil)

	// A recovery code works once.
	phone := s.client()
	phone.loginAnswer(superEmail, superPassword)
	phone.must(http.StatusOK, http.MethodPost, "/auth/mfa", map[string]string{"code": confirmed.RecoveryCodes[0]}, nil)

	again := s.client()
	again.loginAnswer(superEmail, superPassword)
	if status := again.do(http.MethodPost, "/auth/mfa", map[string]string{"code": confirmed.RecoveryCodes[0]}, nil); status != http.StatusUnauthorized {
		t.Errorf("reusing a recovery code = %d, want 401", status)
	}

	var status2 struct {
		MFA struct {
			Enabled           bool `json:"enabled"`
			Required          bool `json:"required"`
			RecoveryCodesLeft int  `json:"recovery_codes_left"`
		} `json:"mfa"`
	}
	laptop.must(http.StatusOK, http.MethodGet, "/mfa", nil, &status2)
	if !status2.MFA.Enabled || !status2.MFA.Required || status2.MFA.RecoveryCodesLeft != 9 {
		t.Errorf("mfa status = %+v", status2.MFA)
	}

	// A super admin resets another administrator's factor: they are signed out
	// and set up again at their next sign-in.
	const staffEmail, staffPassword = "staff@example.com", "staff-password-1"
	laptop.must(http.StatusCreated, http.MethodPost, "/admins", map[string]any{
		"email": staffEmail, "first_name": "Staff", "status": "active",
		"password": staffPassword, "confirm_password": staffPassword,
		"assignments": []map[string]any{{"role_id": laptop.adminRoleID("auditor")}},
	}, nil)

	staff := s.client()
	staff.loginAnswer(staffEmail, staffPassword)
	var staffBegun struct {
		Enrolment struct {
			Secret string `json:"secret"`
		} `json:"enrolment"`
	}
	staff.must(http.StatusOK, http.MethodPost, "/mfa/totp", nil, &staffBegun)
	staff.must(http.StatusOK, http.MethodPost, "/mfa/totp/confirm", map[string]string{"code": codeFor(t, staffBegun.Enrolment.Secret, time.Now())}, nil)
	staff.must(http.StatusOK, http.MethodGet, "/me", nil, nil)

	var admins struct {
		Admins []struct {
			ID         string `json:"id"`
			Email      string `json:"email"`
			MFAEnabled bool   `json:"mfa_enabled"`
		} `json:"admins"`
	}
	laptop.must(http.StatusOK, http.MethodGet, "/admins", nil, &admins)
	staffID := ""
	for _, a := range admins.Admins {
		if a.Email == staffEmail {
			staffID = a.ID
			if !a.MFAEnabled {
				t.Error("the admin list says staff has no second factor")
			}
		}
	}

	var self struct {
		Admin struct {
			ID string `json:"id"`
		} `json:"admin"`
	}
	laptop.must(http.StatusOK, http.MethodGet, "/me", nil, &self)
	if status := laptop.do(http.MethodDelete, "/admins/"+self.Admin.ID+"/mfa", nil, nil); status != http.StatusConflict {
		t.Errorf("resetting your own factor = %d, want 409", status)
	}

	laptop.must(http.StatusOK, http.MethodDelete, "/admins/"+staffID+"/mfa", nil, nil)
	if status := staff.do(http.MethodGet, "/me", nil, nil); status != http.StatusUnauthorized {
		t.Errorf("staff after a reset = %d, want signed out", status)
	}
	if _, answer := staff.loginAnswer(staffEmail, staffPassword); answer.Next != "enroll" {
		t.Errorf("staff login after a reset next = %q, want enroll", answer.Next)
	}
}

// Where it is optional, an administrator signs in with the password alone
// until they turn it on, and can turn it off again with a code.
func TestLiveAdminMFAOptional(t *testing.T) {
	s := newLiveServer(t)
	root := s.superAdmin()

	var begun struct {
		Enrolment struct {
			Secret string `json:"secret"`
		} `json:"enrolment"`
	}
	root.must(http.StatusOK, http.MethodPost, "/mfa/totp", nil, &begun)
	now := time.Now()
	root.must(http.StatusOK, http.MethodPost, "/mfa/totp/confirm", map[string]string{"code": codeFor(t, begun.Enrolment.Secret, now)}, nil)

	// Replacing it takes a code from the current one.
	if status := root.do(http.MethodPost, "/mfa/totp", nil, nil); status != http.StatusConflict {
		t.Errorf("starting a replacement without a code = %d, want 409", status)
	}

	var codes struct {
		RecoveryCodes []string `json:"recovery_codes"`
	}
	root.must(http.StatusOK, http.MethodPost, "/mfa/recovery-codes", map[string]string{"code": codeFor(t, begun.Enrolment.Secret, now.Add(totp.Period))}, &codes)
	if len(codes.RecoveryCodes) != 10 {
		t.Fatalf("regenerated codes = %v", codes.RecoveryCodes)
	}

	root.must(http.StatusNoContent, http.MethodDelete, "/mfa/totp", map[string]string{"code": codes.RecoveryCodes[0]}, nil)

	other := s.client()
	if _, answer := other.loginAnswer(superEmail, superPassword); answer.Admin == nil || answer.Next != "" {
		t.Errorf("login after turning it off = %+v, want signed in with the password", answer)
	}
}
