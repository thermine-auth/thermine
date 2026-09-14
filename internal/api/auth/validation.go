package auth

import (
	"strings"

	"xermess/internal/api/validate"
)

// validate checks the sign-in form before anything is looked up.
func (r *loginRequest) validate() error {
	r.Username = strings.TrimSpace(r.Username)

	return validate.Struct(r)
}
