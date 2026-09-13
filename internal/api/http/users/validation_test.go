package users

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"xermess/internal/api/http/respond"
)

// TestUserRequestValidate covers the one column a user record has of its own.
// It is what someone signs in with, so an address that is not one has no
// business being stored.
func TestUserRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request userRequest
		want    string // the message, or "" when the record is fine
	}{
		{
			name:    "an address",
			request: userRequest{Email: "mira@example.com"},
		},
		{
			name:    "spaces and capitals are tidied away",
			request: userRequest{Email: "  Mira@Example.com "},
		},
		{
			name:    "nothing at all",
			request: userRequest{},
			want:    "email is required",
		},
		{
			name:    "spaces only",
			request: userRequest{Email: "   "},
			want:    "email is required",
		},
		{
			name:    "not an address",
			request: userRequest{Email: "mira at example"},
			want:    "email must be an email address",
		},
		{
			name:    "longer than the column",
			request: userRequest{Email: strings.Repeat("m", 250) + "@example.com"},
			want:    "email must be at most 255 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.request
			err := request.validate()

			if tt.want == "" {
				if err != nil {
					t.Fatalf("validate() = %v, want nothing", err)
				}
				if request.Email != "mira@example.com" {
					t.Errorf("email = %q, want it trimmed and lowered", request.Email)
				}
				return
			}

			var fault respond.Fault
			if !errors.As(err, &fault) {
				t.Fatalf("validate() = %v, want a respond.Fault", err)
			}

			if fault.Status != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", fault.Status)
			}
			if fault.Message != tt.want {
				t.Errorf("message = %q, want %q", fault.Message, tt.want)
			}
		})
	}
}
