package model

import "testing"

func TestRoleNameValid(t *testing.T) {
	valid := []RoleName{RoleSuperAdmin, RoleAdmin, RoleSupport, RoleAuditor}
	for _, role := range valid {
		if !role.Valid() {
			t.Errorf("%s.Valid() = false, want true", role)
		}
	}

	invalid := []RoleName{"", "Admin", "root", "super admin"}
	for _, role := range invalid {
		if role.Valid() {
			t.Errorf("%q.Valid() = true, want false", role)
		}
	}
}

func TestStatusValid(t *testing.T) {
	valid := []Status{StatusInvited, StatusActive, StatusSuspended, StatusDisabled}
	for _, status := range valid {
		if !status.Valid() {
			t.Errorf("%s.Valid() = false, want true", status)
		}
	}

	invalid := []Status{"", "Active", "deleted"}
	for _, status := range invalid {
		if status.Valid() {
			t.Errorf("%q.Valid() = true, want false", status)
		}
	}
}

func TestMFAMethodValid(t *testing.T) {
	if !MFAMethodTOTP.Valid() || !MFAMethodWebAuthn.Valid() {
		t.Error("a known method reported itself invalid")
	}
	if MFAMethod("sms").Valid() {
		t.Error(`MFAMethod("sms").Valid() = true, want false`)
	}
}
