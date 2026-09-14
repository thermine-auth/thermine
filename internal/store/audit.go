package store

import (
	"context"
	"time"

	"xermess/internal/model"
)

// WriteAudit adds one line to the activity log.
func (s *Store) WriteAudit(ctx context.Context, entry *model.AuditLog) error {
	return s.db.WithContext(ctx).Create(entry).Error
}

// AuditLog returns the newest entries, which is what both the dashboard and
// the logs page show.
func (s *Store) AuditLog(ctx context.Context, limit int) ([]model.AuditLog, error) {
	var events []model.AuditLog
	err := s.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&events).Error

	return events, err
}

// Counts are the numbers on the dashboard.
type Counts struct {
	Admins         int64 `json:"admins"`
	ActiveSessions int64 `json:"active_sessions"`
	Applications   int64 `json:"applications"`
	Roles          int64 `json:"roles"`
	Events         int64 `json:"events"`
}

// Counts totals the tables the dashboard reports on. A count that fails takes
// the whole answer with it: a dashboard of partly wrong numbers is worse than
// an error.
func (s *Store) Counts(ctx context.Context, now time.Time) (Counts, error) {
	db := s.db.WithContext(ctx)
	var counts Counts

	if err := db.Model(&model.AdminUser{}).Count(&counts.Admins).Error; err != nil {
		return counts, err
	}

	err := db.Model(&model.AdminUserSession{}).
		Where("revoked_at IS NULL AND expires_at > ?", now).
		Count(&counts.ActiveSessions).Error
	if err != nil {
		return counts, err
	}

	if err := db.Model(&model.Application{}).Count(&counts.Applications).Error; err != nil {
		return counts, err
	}

	if err := db.Model(&model.Role{}).Count(&counts.Roles).Error; err != nil {
		return counts, err
	}

	if err := db.Model(&model.AuditLog{}).Count(&counts.Events).Error; err != nil {
		return counts, err
	}

	return counts, nil
}
