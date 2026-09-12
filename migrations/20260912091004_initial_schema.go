package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"xermess/internal/model"
)

func init() {
	goose.AddMigrationContext(upInitialSchema, downInitialSchema)
}

// upInitialSchema creates every table the models describe: admin_users,
// admin_user_sessions, roles, permissions, mfa, audit_logs, and the two join
// tables that carry the roles and their permissions.
func upInitialSchema(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	return db.AutoMigrate(model.All()...)
}

// downInitialSchema drops those tables again, children before parents so the
// foreign keys do not block it.
func downInitialSchema(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	return db.Migrator().DropTable(
		"role_permissions",
		"admin_user_roles",
		&model.MFA{},
		&model.AdminUserSession{},
		&model.AuditLog{},
		&model.Permission{},
		&model.Role{},
		&model.AdminUser{},
	)
}
