package account

import (
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"xermess/internal/api/respond"
	"xermess/internal/api/validate"
	"xermess/internal/oidc"
)

func (r *loginRequest) validate() error {
	r.Email = strings.TrimSpace(r.Email)
	return validate.Struct(r)
}

func (r *registerRequest) validate() error {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)

	if err := validate.Struct(r); err != nil {
		return err
	}

	return checkPassword(r.Password)
}

func (r *forgotRequest) validate() error {
	r.Email = strings.TrimSpace(r.Email)
	return validate.Struct(r)
}

func (r *resetRequest) validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}

	return checkPassword(r.Password)
}

func (r *profileRequest) validate() error {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	return validate.Struct(r)
}

func (r *passwordRequest) validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}

	return checkPassword(r.NewPassword)
}

// checkPassword holds a new password to the one rule there is: its length.
// bcrypt reads no further than 72 bytes, so a longer one is refused rather
// than quietly cut short.
func checkPassword(password string) error {
	switch {
	case utf8.RuneCountInString(password) < oidc.MinPasswordLength:
		return respond.Fault{Status: http.StatusBadRequest, Message: fmt.Sprintf("password must be at least %d characters", oidc.MinPasswordLength)}
	case len(password) > 72:
		return respond.Fault{Status: http.StatusBadRequest, Message: "password must be at most 72 bytes"}
	}

	return nil
}
