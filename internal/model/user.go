package model

// User is an account that signs in to the applications this server protects,
// as opposed to AdminUser, who signs in to the panel.
//
// Only the fields every user has are columns. Everything else an organisation
// wants to keep — a first name, whether a phone was verified — is defined in
// user_fields and stored in Data, so adding one does not need a migration.
type User struct {
	Base

	Email         string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	EmailVerified bool   `gorm:"not null;default:false" json:"email_verified"`

	Data map[string]any `gorm:"serializer:json" json:"data"`
}

// TableName pins the table name.
func (User) TableName() string {
	return "users"
}
