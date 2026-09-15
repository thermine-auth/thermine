package oidc

import (
	"errors"
	"net/http"
)

// The error codes of RFC 6749 section 4.1.2.1 and 5.2, RFC 6750 section 3.1
// and OpenID Connect Core section 3.1.2.6 that this server answers with.
const (
	ErrInvalidRequest          = "invalid_request"
	ErrInvalidClient           = "invalid_client"
	ErrInvalidGrant            = "invalid_grant"
	ErrUnauthorizedClient      = "unauthorized_client"
	ErrUnsupportedGrantType    = "unsupported_grant_type"
	ErrUnsupportedResponseType = "unsupported_response_type"
	ErrInvalidScope            = "invalid_scope"
	ErrAccessDenied            = "access_denied"
	ErrServerError             = "server_error"
	ErrLoginRequired           = "login_required"
	ErrInvalidToken            = "invalid_token"
	ErrInsufficientScope       = "insufficient_scope"
)

// Error is an OAuth error: a code a client can act on, a sentence a developer
// can read, and the HTTP status the endpoints answer it with.
type Error struct {
	Code        string
	Description string
	Status      int
}

func (e *Error) Error() string {
	return e.Code + ": " + e.Description
}

func oauthError(code, description string) *Error {
	status := http.StatusBadRequest
	switch code {
	case ErrInvalidClient, ErrInvalidToken:
		status = http.StatusUnauthorized
	case ErrInsufficientScope:
		status = http.StatusForbidden
	}

	return &Error{Code: code, Description: description, Status: status}
}

// The errors of signing a user in. They are what the sign-in pages show, so
// they say what went wrong without saying which accounts exist.
var (
	// ErrInvalidCredentials covers an unknown address, a wrong password, and
	// an account that may not sign in.
	ErrInvalidCredentials = errors.New("wrong email or password")
	// ErrRequestExpired is a sign-in handle that has expired, been used, or
	// never existed.
	ErrRequestExpired = errors.New("this sign-in link has expired; go back to the application and start again")
	// ErrRegistrationClosed is registering where the application does not
	// allow it.
	ErrRegistrationClosed = errors.New("this application does not allow creating an account")
	// ErrEmailTaken is registering an address that already has an account.
	ErrEmailTaken = errors.New("an account with this email already exists")
	// ErrResetInvalid is a reset link that has expired or been used.
	ErrResetInvalid = errors.New("this reset link has expired or has already been used")
)

// FieldError is a problem with what someone typed into a sign-in form, which
// the page shows as it is.
type FieldError struct {
	Message string
}

func (e *FieldError) Error() string {
	return e.Message
}
