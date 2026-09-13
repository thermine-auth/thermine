package model

import "time"

// AdminUser is a member of staff who can sign in to the admin panel.
type AdminUser struct {
	Base

	Username     string `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email        string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	FirstName    string `gorm:"size:100;not null" json:"first_name"`
	LastName     string `gorm:"size:100;not null" json:"last_name"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	IsActive     bool   `gorm:"-" json:"is_active"`
	Status       Status `gorm:"type:varchar(32);index;not null;default:invited" json:"status"`

	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP     string     `gorm:"size:45" json:"last_login_ip,omitempty"`

	// FailedLoginCount and LockedUntil back a lockout after repeated failures.
	FailedLoginCount int        `gorm:"not null;default:0" json:"-"`
	LockedUntil      *time.Time `json:"locked_until,omitempty"`

	Roles    []Role             `gorm:"many2many:admin_user_roles;constraint:OnDelete:CASCADE" json:"roles,omitempty"`
	Sessions []AdminUserSession `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	MFA      []MFA              `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

// TableName pins the table name so renaming the struct cannot silently rename
// the table.
func (AdminUser) TableName() string {
	return "admin_users"
}

// FullName is the admin's display name.
func (a AdminUser) FullName() string {
	return a.FirstName + " " + a.LastName
}

// CanSignIn reports whether the account is in a state that allows signing in.
func (a AdminUser) CanSignIn(now time.Time) bool {
	if a.Status != StatusActive {
		return false
	}
	return a.LockedUntil == nil || now.After(*a.LockedUntil)
}

// HasRole reports whether the admin holds the given role. Roles must be loaded
// for this to mean anything.
func (a AdminUser) HasRole(name RoleName) bool {
	for _, r := range a.Roles {
		if r.Name == name {
			return true
		}
	}
	return false
}

// HasPermission reports whether any of the admin's roles grant the permission.
// Roles and their permissions must be loaded for this to mean anything.
func (a AdminUser) HasPermission(name string) bool {
	for _, r := range a.Roles {
		if r.Name == RoleSuperAdmin {
			return true
		}
		for _, p := range r.Permissions {
			if p.Name == name {
				return true
			}
		}
	}
	return false
}
