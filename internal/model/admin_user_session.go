package model

import (
	"time"

	"github.com/google/uuid"
)

// AdminUserSession is one issued refresh token. Only the hash of the token is
// stored, so a leaked database cannot be replayed as a login.
type AdminUserSession struct {
	Base

	AdminUserID uuid.UUID `gorm:"type:uuid;index;not null" json:"admin_user_id"`
	AdminUser   AdminUser `gorm:"foreignKey:AdminUserID" json:"-"`

	TokenHash string     `gorm:"uniqueIndex;size:64;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"index;not null" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`

	// MFAPassed records whether the second factor was cleared for this
	// session, so a half-finished login cannot be used.
	MFAPassed bool `gorm:"not null;default:false" json:"mfa_passed"`

	UserAgent string `gorm:"size:255" json:"user_agent,omitempty"`
	IP        string `gorm:"size:45" json:"ip,omitempty"`
}

// TableName pins the table name.
func (AdminUserSession) TableName() string {
	return "admin_user_sessions"
}

// IsActive reports whether the session can still be exchanged for a new token.
func (s AdminUserSession) IsActive(now time.Time) bool {
	return s.RevokedAt == nil && s.MFAPassed && now.Before(s.ExpiresAt)
}
