package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"xermess/internal/model"
)

// UserQuery is what a listing asks for: a search box, a filter, and a page.
type UserQuery struct {
	// Search matches the email or any value the record holds.
	Search string
	// Verified filters on the email_verified flag. Nil means both.
	Verified *bool
	Limit    int
	Offset   int
}

// Users returns a page of users, newest first, along with how many match the
// query in total — which is what the panel shows above the table.
func (s *Store) Users(ctx context.Context, q UserQuery) ([]model.User, int64, error) {
	query := s.db.WithContext(ctx).Model(&model.User{})

	if search := strings.TrimSpace(q.Search); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		// The fields live in a JSON column; casting it to text searches every
		// one of them at once, which is what a single search box should do.
		query = query.Where("LOWER(email) LIKE ? OR LOWER(data::text) LIKE ?", like, like)
	}

	if q.Verified != nil {
		query = query.Where("email_verified = ?", *q.Verified)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.User
	err := query.Order("created_at DESC").Limit(q.Limit).Offset(q.Offset).Find(&users).Error

	return users, total, err
}

// User returns one user by id.
func (s *Store) User(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, translate(err)
	}

	return &user, nil
}

// CreateUser writes a new user.
func (s *Store) CreateUser(ctx context.Context, user *model.User) error {
	return translate(s.db.WithContext(ctx).Create(user).Error)
}

// SaveUser writes a user back, columns and fields alike.
func (s *Store) SaveUser(ctx context.Context, user *model.User) error {
	return translate(s.db.WithContext(ctx).Save(user).Error)
}

// DeleteUser removes a user for good rather than marking it deleted: an
// account someone asked to have removed should not stay in the table.
func (s *Store) DeleteUser(ctx context.Context, user *model.User) error {
	return s.db.WithContext(ctx).Unscoped().Delete(user).Error
}

// FieldValueTaken reports whether another user already holds this value for a
// field that has to be unique. `except` is the record being written, so it
// does not clash with itself; pass uuid.Nil when creating one.
//
// The values live in a JSON column, so this is a query rather than an index.
func (s *Store) FieldValueTaken(ctx context.Context, field string, value any, except uuid.UUID) (bool, error) {
	text := fmt.Sprintf("%v", value)
	if number, ok := value.(float64); ok {
		text = strconv.FormatFloat(number, 'f', -1, 64)
	}

	query := s.db.WithContext(ctx).
		Model(&model.User{}).
		Where("data::jsonb ->> ? = ?", field, text)

	if except != uuid.Nil {
		query = query.Where("id <> ?", except)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
