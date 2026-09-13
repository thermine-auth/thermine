package store

import (
	"context"
	"time"

	"github.com/google/uuid"

	"xermess/internal/model"
)

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
