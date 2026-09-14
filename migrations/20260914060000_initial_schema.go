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

// upInitialSchema builds the whole database in one step: every table the
// models describe, including the tables joining users to their roles and
// roles to the roles they inherit, and the admin roles a panel starts with.
//
// It seeds no applications, user fields or user roles, global or otherwise. The fields every record has — an
// address, a name, a password, whether the account is active — are columns
// of the users table, created here from the model; everything else, and
// every role users hold, is added in the panel, so a new installation starts
// with none.
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

	// A null application means the whole panel for an admin role, and a
	// global role for a user role. Postgres treats every null as different,
	// so a unique index over a nullable column needs a partial index for the
	// null case.
	indexes := []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_user_roles_global_name
			ON user_roles (name) WHERE application_id IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_role_assignments_global
			ON admin_role_assignments (admin_user_id, role_id) WHERE application_id IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_role_assignments_scoped
			ON admin_role_assignments (admin_user_id, role_id, application_id) WHERE application_id IS NOT NULL`,
	}
	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			return fmt.Errorf("create assignment index: %w", err)
		}
	}

	return seedAdminRoles(db)
}

// downInitialSchema drops everything again, join tables and children before
// the tables they point at, so the foreign keys do not block it.
func downInitialSchema(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	return db.Migrator().DropTable(
		&model.AdminRoleAssignment{},
		"user_role_members",
		"user_role_inherits",
		"user_role_api_scopes",
		&model.ApplicationAPIScope{},
		&model.ApplicationAPI{},
		&model.APIScope{},
		&model.API{},
		&model.MFA{},
		&model.AdminUserSession{},
		&model.AuditLog{},
		&model.AdminUser{},
		&model.Role{},
		&model.User{},
		&model.UserRole{},
		&model.Application{},
		&model.UserField{},
	)
}

// startingRoles are the admin roles a panel is given, and what each grants.
//
// super_admin is built in: it grants every permission by name, so it lists
// none, and it is the only role that can manage administrators. A super admin
// can change or remove any of the others. app_manager is meant to be held for
// one application: every one of its permissions can be scoped.
var startingRoles = []model.Role{
	{
		Name:        model.RoleSuperAdmin,
		Description: "Everything, including managing administrators",
	},
	{
		Name:        "admin",
		Description: "Everything but managing administrators",
		Permissions: model.AdminPermissionNames(),
	},
	{
		Name:        "moderator",
		Description: "Manage users, application roles and who holds them, and read the activity log",
		Permissions: []string{
			model.PermActivityRead, model.PermUsersRead, model.PermUsersWrite, model.PermAPIsRead,
			model.PermApplicationsRead, model.PermUserRolesWrite, model.PermRoleAssignmentsWrite,
		},
	},
	{
		Name:        "app_manager",
		Description: "Look after an application: its settings, its roles, and who holds them",
		Permissions: []string{
			model.PermApplicationsRead, model.PermApplicationsWrite,
			model.PermUserRolesWrite, model.PermRoleAssignmentsWrite,
		},
	},
	{
		Name:        "support",
		Description: "Look users up, fix their records and give them application roles",
		Permissions: []string{
			model.PermUsersRead, model.PermUsersWrite, model.PermApplicationsRead, model.PermRoleAssignmentsWrite,
		},
	},
	{
		Name:        "auditor",
		Description: "Read only, including the activity log",
		Permissions: []string{
			model.PermActivityRead, model.PermUsersRead, model.PermAPIsRead, model.PermApplicationsRead,
		},
	},
}

// seedAdminRoles creates the starting roles. A role that is already there is
// left alone, so this is safe to run twice.
func seedAdminRoles(db *gorm.DB) error {
	for _, starting := range startingRoles {
		role := starting

		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error; err != nil {
			return fmt.Errorf("create role %s: %w", starting.Name, err)
		}
	}

	return nil
}
