package model

import (
	"testing"
	"time"
)

func TestFullName(t *testing.T) {
	admin := AdminUser{FirstName: "Ada", LastName: "Lovelace"}

	if got, want := admin.FullName(), "Ada Lovelace"; got != want {
		t.Errorf("FullName() = %q, want %q", got, want)
	}
}

func TestCanSignIn(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name  string
		admin AdminUser
		want  bool
	}{
		{
			name:  "active account",
			admin: AdminUser{Status: StatusActive},
			want:  true,
		},
		{
			name:  "invited account has not set a password yet",
			admin: AdminUser{Status: StatusInvited},
		},
		{
			name:  "suspended account",
			admin: AdminUser{Status: StatusSuspended},
		},
		{
			name:  "disabled account",
			admin: AdminUser{Status: StatusDisabled},
		},
		{
			name:  "active but locked out after failed attempts",
			admin: AdminUser{Status: StatusActive, LockedUntil: &future},
		},
		{
			name:  "active and the lockout has expired",
			admin: AdminUser{Status: StatusActive, LockedUntil: &past},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.admin.CanSignIn(now); got != tt.want {
				t.Errorf("CanSignIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasRole(t *testing.T) {
	admin := AdminUser{Roles: []Role{{Name: RoleSupport}}}

	if !admin.HasRole(RoleSupport) {
		t.Error("HasRole(support) = false, want true")
	}
	if admin.HasRole(RoleAdmin) {
		t.Error("HasRole(admin) = true, want false")
	}

	// Roles have to be loaded for this to mean anything, so an admin with
	// none reports false rather than panicking.
	if (AdminUser{}).HasRole(RoleAdmin) {
		t.Error("an admin with no roles loaded reported having one")
	}
}

func TestHasPermission(t *testing.T) {
	read := Permission{Name: "admin_users.read"}
	write := Permission{Name: "admin_users.create"}

	tests := []struct {
		name       string
		admin      AdminUser
		permission string
		want       bool
	}{
		{
			name:       "role grants the permission",
			admin:      AdminUser{Roles: []Role{{Name: RoleSupport, Permissions: []Permission{read}}}},
			permission: "admin_users.read",
			want:       true,
		},
		{
			name:       "role does not grant the permission",
			admin:      AdminUser{Roles: []Role{{Name: RoleSupport, Permissions: []Permission{read}}}},
			permission: "admin_users.create",
		},
		{
			name:       "permission comes from the second of two roles",
			admin:      AdminUser{Roles: []Role{{Name: RoleAuditor, Permissions: []Permission{read}}, {Name: RoleSupport, Permissions: []Permission{write}}}},
			permission: "admin_users.create",
			want:       true,
		},
		{
			name:       "super admin may do anything, with no permissions loaded",
			admin:      AdminUser{Roles: []Role{{Name: RoleSuperAdmin}}},
			permission: "anything.at.all",
			want:       true,
		},
		{
			name:       "no roles at all",
			admin:      AdminUser{},
			permission: "admin_users.read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.admin.HasPermission(tt.permission); got != tt.want {
				t.Errorf("HasPermission(%q) = %v, want %v", tt.permission, got, tt.want)
			}
		})
	}
}

func TestSessionIsActive(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name    string
		session AdminUserSession
		want    bool
	}{
		{
			name:    "current session that cleared MFA",
			session: AdminUserSession{ExpiresAt: future, MFAPassed: true},
			want:    true,
		},
		{
			name:    "expired",
			session: AdminUserSession{ExpiresAt: past, MFAPassed: true},
		},
		{
			name:    "revoked",
			session: AdminUserSession{ExpiresAt: future, MFAPassed: true, RevokedAt: &past},
		},
		{
			name:    "half-finished login that never cleared MFA",
			session: AdminUserSession{ExpiresAt: future},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.session.IsActive(now); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMFAIsConfirmed(t *testing.T) {
	now := time.Now()

	if (MFA{}).IsConfirmed() {
		t.Error("an unconfirmed factor reported itself confirmed")
	}
	if !(MFA{ConfirmedAt: &now}).IsConfirmed() {
		t.Error("a confirmed factor reported itself unconfirmed")
	}
}
