package model

// Permission is a single thing an admin may do, named "<resource>.<action>",
// for example "admin_users.create".
type Permission struct {
	Base

	Name        string `gorm:"uniqueIndex;size:128;not null" json:"name"`
	Resource    string `gorm:"index;size:64;not null" json:"resource"`
	Action      string `gorm:"size:64;not null" json:"action"`
	Description string `gorm:"size:255" json:"description,omitempty"`

	Roles []Role `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName pins the table name.
func (Permission) TableName() string {
	return "permissions"
}
