package model

type User struct {
	Base

	Email         string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	EmailVerified bool   `gorm:"not null;default:false" json:"email_verified"`
	FirstName     string `gorm:"size:100" json:"first_name"`
	LastName      string `gorm:"size:100" json:"last_name"`
	IsActive      bool   `gorm:"not null;default:true" json:"is_active"`

	Data map[string]any `gorm:"serializer:json" json:"data"`
}

func (User) TableName() string {
	return "users"
}

func (u User) FullName() string {
	switch {
	case u.FirstName != "" && u.LastName != "":
		return u.FirstName + " " + u.LastName
	case u.FirstName != "":
		return u.FirstName
	case u.LastName != "":
		return u.LastName
	default:
		return u.Email
	}
}
