package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"xermess/internal/model"
)

func init() {
	goose.AddMigrationContext(upInitialSchema, downInitialSchema)
}

// upInitialSchema builds the whole database: every table the models describe,
// and the roles an administrator can hold.
//
// It seeds no user fields. The fields every record has — an address, a name,
// whether the account is active — are columns of the users table, created
// here from the model; everything else, a phone number included, is added by
// a super admin in the panel, so a new installation starts with none.
//
// It creates no administrator either. The first one is made by whoever opens
// the panel, on /admin/new-super-admin, which is offered only while there is
// none — so a deployment has no password written down anywhere, and no
// default one to forget to change.
func upInitialSchema(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(model.All()...); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	return seedRoles(db)
}

// downInitialSchema drops everything again, children before parents so the
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
		&model.User{},
		&model.UserField{},
	)
}

// seedRoles creates the roles the models define. Anything already there is
// left alone, so this is safe to run twice.
func seedRoles(db *gorm.DB) error {
	for _, name := range model.RoleNames {
		role := model.Role{Name: name, Description: string(name) + " role"}

		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error; err != nil {
			return fmt.Errorf("create role %s: %w", name, err)
		}
	}

	return nil
}
