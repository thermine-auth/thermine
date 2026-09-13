package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"xermess/internal/model"
)

// ErrAdminExists is returned when the first administrator is asked for and
// there already is one. It is what keeps the endpoint that creates it from
// being a way in once the panel is set up.
var ErrAdminExists = errors.New("an administrator already exists")

// AdminsExist reports whether anyone can sign in to the panel yet.
func (s *Store) AdminsExist(ctx context.Context) (bool, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.AdminUser{}).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateFirstAdmin writes the administrator a new installation is set up
// with, and refuses if there is one already.
//
// The check and the write are one transaction, so two people opening the
// setup page at the same moment cannot both get an account: the second finds
// the first and is turned away.
func (s *Store) CreateFirstAdmin(ctx context.Context, admin *model.AdminUser) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.AdminUser{}).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return ErrAdminExists
		}

		var role model.Role
		if err := tx.Where("name = ?", model.RoleSuperAdmin).First(&role).Error; err != nil {
			return err
		}

		admin.Roles = []model.Role{role}

		return translate(tx.Create(admin).Error)
	})
}

// AdminByUsername loads an administrator and the roles they hold. It is what
// signing in starts with.
func (s *Store) AdminByUsername(ctx context.Context, username string) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := s.db.WithContext(ctx).
		Preload("Roles").
		Where("username = ?", username).
		First(&admin).Error
	if err != nil {
		return nil, translate(err)
	}

	return &admin, nil
}

// AdminByID loads an administrator with their roles and everything those
// roles allow, which is what a request is checked against.
func (s *Store) AdminByID(ctx context.Context, id uuid.UUID) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := s.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(&admin, "id = ?", id).Error
	if err != nil {
		return nil, translate(err)
	}

	return &admin, nil
}

// MarkAdminSignedIn records when and from where an administrator last signed
// in.
func (s *Store) MarkAdminSignedIn(ctx context.Context, admin *model.AdminUser, at time.Time, ip string) error {
	return s.db.WithContext(ctx).Model(admin).
		Updates(map[string]any{"last_login_at": at, "last_login_ip": ip}).Error
}
