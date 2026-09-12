package model

// RoleName is the name of a role. Roles live in their own table so they can
// carry permissions, but the names an admin can hold are fixed here.
type RoleName string

const (
	RoleSuperAdmin RoleName = "super_admin"
	RoleAdmin      RoleName = "admin"
	RoleSupport    RoleName = "support"
	RoleAuditor    RoleName = "auditor"
)

// RoleNames lists every role, for seeding and validation.
var RoleNames = []RoleName{RoleSuperAdmin, RoleAdmin, RoleSupport, RoleAuditor}

// Valid reports whether r is a known role.
func (r RoleName) Valid() bool {
	for _, known := range RoleNames {
		if r == known {
			return true
		}
	}
	return false
}

// String returns the role name.
func (r RoleName) String() string {
	return string(r)
}

// Status is the state of an admin account.
type Status string

const (
	// StatusInvited is an account that has been created but has not yet set a
	// password.
	StatusInvited Status = "invited"
	// StatusActive is an account that can sign in.
	StatusActive Status = "active"
	// StatusSuspended is an account temporarily blocked from signing in.
	StatusSuspended Status = "suspended"
	// StatusDisabled is an account permanently blocked from signing in.
	StatusDisabled Status = "disabled"
)

// Statuses lists every status, for validation.
var Statuses = []Status{StatusInvited, StatusActive, StatusSuspended, StatusDisabled}

// Valid reports whether s is a known status.
func (s Status) Valid() bool {
	for _, known := range Statuses {
		if s == known {
			return true
		}
	}
	return false
}

// String returns the status.
func (s Status) String() string {
	return string(s)
}

// MFAMethod is the kind of second factor an admin has enrolled.
type MFAMethod string

const (
	MFAMethodTOTP     MFAMethod = "totp"
	MFAMethodWebAuthn MFAMethod = "webauthn"
)

// MFAMethods lists every method, for validation.
var MFAMethods = []MFAMethod{MFAMethodTOTP, MFAMethodWebAuthn}

// Valid reports whether m is a known method.
func (m MFAMethod) Valid() bool {
	for _, known := range MFAMethods {
		if m == known {
			return true
		}
	}
	return false
}

// String returns the method.
func (m MFAMethod) String() string {
	return string(m)
}
