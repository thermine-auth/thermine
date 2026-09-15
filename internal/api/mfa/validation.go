package mfa

// The rules a code keeps — six digits for an authenticator, ten characters for
// a recovery code — are checked where the code is checked, in internal/auth
// and internal/totp: a code in the wrong shape is simply a wrong code, and says
// so the same way.
