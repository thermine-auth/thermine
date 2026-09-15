package account

import (
	"time"

	"github.com/google/uuid"

	"xermess/internal/model"
	"xermess/internal/oidc"
)

// requestResponse is a sign-in under way, as the sign-in page shows it.
type requestResponse struct {
	Application oidc.PublicApplication `json:"application"`
	LoginHint   string                 `json:"login_hint"`
	ExpiresAt   time.Time              `json:"expires_at"`
}

func newRequestResponse(pending *oidc.PendingRequest) requestResponse {
	return requestResponse{
		Application: oidc.Public(pending.Application),
		LoginHint:   pending.LoginHint,
		ExpiresAt:   pending.ExpiresAt,
	}
}

// signedInResponse says what the page does next: follow RedirectTo back to
// the application, or — for a temporary password — send the user to choose a
// new one with ResetToken. With neither, the user is signed in and there is
// nowhere to go.
type signedInResponse struct {
	RedirectTo             string `json:"redirect_to,omitempty"`
	PasswordChangeRequired bool   `json:"password_change_required,omitempty"`
	ResetToken             string `json:"reset_token,omitempty"`
}

// userResponse is a user as they see themselves. It is built by hand, so a
// column added to users later is not published to the user by accident — the
// roles and the additional fields an organisation keeps are not theirs to read
// here.
type userResponse struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	CreatedAt     time.Time  `json:"created_at"`
	LastLoginAt   *time.Time `json:"last_login_at"`
}

func newUserResponse(user *model.User) userResponse {
	return userResponse{
		ID:            user.ID,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		CreatedAt:     user.CreatedAt,
		LastLoginAt:   user.LastLoginAt,
	}
}
