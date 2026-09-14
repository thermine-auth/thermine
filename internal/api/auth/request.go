package auth

// loginRequest is the body of POST /admin/auth/login.
//
// The password has no rule beyond being there: what a password must look like
// is the account's business, and a sign-in form that explains the policy is a
// sign-in form that helps someone guessing.
type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
