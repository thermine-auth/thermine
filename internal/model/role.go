package model

// Role is a named bundle of permissions. Admins hold roles; roles hold
// permissions. The set of role names is fixed in RoleNames.
type Role struct {
	Base

	Name        RoleName `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`
	Description string   `gorm:"size:255" json:"description,omitempty"`

	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE" json:"permissions,omitempty"`
	AdminUsers  []AdminUser  `gorm:"many2many:admin_user_roles;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName pins the table name.
func (Role) TableName() string {
	return "roles"
}
