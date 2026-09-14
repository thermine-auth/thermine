package model

import (
	"regexp"
	"slices"
	"sort"

	"github.com/google/uuid"
)

// UserRole is a role users hold, either across every application or in one.
//
// A global role belongs to no application — Keycloak's realm roles — and is
// the same role wherever the user signs in: "employee", "beta-tester". An
// application role belongs to the application that defines it — Keycloak's
// client roles, Zitadel's project roles, Entra's app roles — so "admin" in
// one app says nothing about another. A name is unique within its scope: once
// among the global roles, and once within each application.
//
// This server says which roles a user holds; what a role lets them do is
// decided by the applications, so a role carries no permissions of its own.
//
// A role can be composite, inheriting other roles: holding "editor" also
// means holding "viewer". As in Keycloak, a global role may include global
// roles and any application's roles, and an application role may include
// global roles and roles of its own application. A role marked default is
// given to every user created.
type UserRole struct {
	Base

	// ApplicationID is the application the role belongs to, or nil for a
	// global role.
	ApplicationID *uuid.UUID   `gorm:"type:uuid;uniqueIndex:idx_user_roles_application_name,priority:1" json:"application_id"`
	Application   *Application `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	Name        string `gorm:"size:64;not null;uniqueIndex:idx_user_roles_application_name,priority:2" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	IsDefault   bool   `gorm:"not null;default:false;index" json:"is_default"`

	Inherits []UserRole `gorm:"many2many:user_role_inherits;joinForeignKey:RoleID;joinReferences:InheritedRoleID;constraint:OnDelete:CASCADE" json:"inherits,omitempty"`

	// APIScopes are the API scopes holding this role grants. A token for an
	// API carries one only when the application is allowed it too.
	APIScopes []APIScope `gorm:"many2many:user_role_api_scopes;constraint:OnDelete:CASCADE" json:"api_scopes,omitempty"`
}

// Global reports whether the role belongs to no application.
func (r UserRole) Global() bool {
	return r.ApplicationID == nil
}

// In reports whether the role belongs to this application. A nil application
// asks about global roles.
func (r UserRole) In(application *uuid.UUID) bool {
	if r.ApplicationID == nil || application == nil {
		return r.ApplicationID == nil && application == nil
	}

	return *r.ApplicationID == *application
}

// MayInherit reports whether this role may include `other`: a global role may
// include any role, and an application role may include global roles and the
// roles of its own application.
func (r UserRole) MayInherit(other UserRole) bool {
	return r.Global() || other.Global() || other.In(r.ApplicationID)
}

// TableName pins the table name, and keeps these apart from the roles
// administrators hold.
func (UserRole) TableName() string {
	return "user_roles"
}

// RoleNamePattern is what a role may be called: short, lower case, and safe
// to put in a token or compare in code. Admin roles keep the same rule.
var RoleNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)

// RoleGraph is every role by id, each with the ids of the roles it inherits
// loaded. Roles are few, so resolving against the whole graph in memory is
// simpler than walking it in SQL.
type RoleGraph map[uuid.UUID]UserRole

// Reaches reports whether `to` can be reached from `from` by following what
// roles inherit, `from` itself included.
func (g RoleGraph) Reaches(from, to uuid.UUID) bool {
	for _, role := range g.walk([]uuid.UUID{from}) {
		if role.ID == to {
			return true
		}
	}

	return false
}

// WouldCycle reports whether letting `role` inherit `inherits` would make a
// role inherit itself, directly or through others. The graph is the one
// stored before the change: a path back to the role through its old
// inheritance still passes through the role, so it is found either way.
func (g RoleGraph) WouldCycle(role uuid.UUID, inherits []uuid.UUID) bool {
	for _, id := range inherits {
		if id == role || g.Reaches(id, role) {
			return true
		}
	}

	return false
}

// Effective is every role holding these amounts to, the held ones included,
// each once and sorted by name. Ids not in the graph are skipped.
func (g RoleGraph) Effective(held []uuid.UUID) []UserRole {
	roles := g.walk(held)
	sort.Slice(roles, func(i, j int) bool { return roles[i].Name < roles[j].Name })

	return roles
}

// walk returns every role reached from these ids, each once, in no particular
// order. A cycle stored by some other route does not make it loop.
func (g RoleGraph) walk(from []uuid.UUID) []UserRole {
	seen := map[uuid.UUID]bool{}
	stack := append([]uuid.UUID(nil), from...)

	var roles []UserRole
	for len(stack) > 0 {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		role, ok := g[id]
		if !ok || seen[id] {
			continue
		}
		seen[id] = true

		roles = append(roles, role)
		for _, inherited := range role.Inherits {
			stack = append(stack, inherited.ID)
		}
	}

	return roles
}

// Via says, for every role holding these amounts to, which of the held roles
// it comes through. A held role counts as coming through itself only when it
// is not also reached through another held role, so a role listed under
// itself is exactly one that was given directly and nothing else implies.
func (g RoleGraph) Via(held []uuid.UUID) map[uuid.UUID][]uuid.UUID {
	via := map[uuid.UUID][]uuid.UUID{}

	for _, start := range held {
		for _, reached := range g.walk([]uuid.UUID{start}) {
			if reached.ID == start {
				continue
			}
			if !slices.Contains(via[reached.ID], start) {
				via[reached.ID] = append(via[reached.ID], start)
			}
		}
	}

	return via
}

// Names returns the roles' names, in the order given.
func Names(roles []UserRole) []string {
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		names = append(names, role.Name)
	}

	return names
}
