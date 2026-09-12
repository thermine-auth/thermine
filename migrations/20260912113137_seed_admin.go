package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"xermess/internal/config"
	"xermess/internal/model"
)

func init() {
	goose.AddMigrationContext(upSeedAdmin, downSeedAdmin)
}

// upSeedAdmin creates the roles and the first administrator, so a fresh
// database can be signed in to. The credentials come from
// XERMESS_ADMIN_USERNAME and XERMESS_ADMIN_PASSWORD; with either unset, the
// roles are still created and no administrator is.
func upSeedAdmin(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	// A database created by the first migration already has this column,
	// because that migration builds the schema from the models, which now
	// carry Username. One created before the field existed does not.
	if !db.Migrator().HasColumn(&model.AdminUser{}, "Username") {
		if err := db.Migrator().AddColumn(&model.AdminUser{}, "Username"); err != nil {
			return fmt.Errorf("add username column: %w", err)
		}
	}

	if err := seedRoles(db); err != nil {
		return err
	}

	admin, err := config.LoadAdmin()
	if err != nil {
		return fmt.Errorf("read configuration: %w", err)
	}

	if admin.Username == "" || admin.Password == "" {
		// Not an error: a test or a review environment has no credentials to
		// seed, and should still end up with a working schema.
		fmt.Println("goose: XERMESS_ADMIN_USERNAME or XERMESS_ADMIN_PASSWORD is empty, no administrator created")
		return nil
	}

	return seedAdmin(db, admin)
}

// downSeedAdmin removes the administrator and the roles.
//
// It leaves the username column alone: on a database built by the first
// migration that column came from the models, and dropping it here would take
// away something this migration did not add.
func downSeedAdmin(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	admin, err := config.LoadAdmin()
	if err != nil {
		return fmt.Errorf("read configuration: %w", err)
	}

	if admin.Username != "" {
		// Unscoped, or the row would only be marked deleted and the unique
		// username would stay taken.
		if err := db.Unscoped().
			Where("username = ?", admin.Username).
			Delete(&model.AdminUser{}).Error; err != nil {
			return err
		}
	}

	return db.Unscoped().
		Where("name IN ?", model.RoleNames).
		Delete(&model.Role{}).Error
}

// seedRoles creates the roles the models define, skipping any that are
// already there, so this is safe to run against a database that has some.
func seedRoles(db *gorm.DB) error {
	for _, name := range model.RoleNames {
		role := model.Role{Name: name, Description: string(name) + " role"}

		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error; err != nil {
			return fmt.Errorf("create role %s: %w", name, err)
		}
	}

	return nil
}

// seedAdmin creates the first administrator, with the super_admin role, and
// does nothing if that username is already taken.
func seedAdmin(db *gorm.DB, admin config.Admin) error {
	var existing int64
	if err := db.Model(&model.AdminUser{}).
		Where("username = ?", admin.Username).
		Count(&existing).Error; err != nil {
		return err
	}
	if existing > 0 {
		fmt.Printf("goose: administrator %q already exists, leaving it alone\n", admin.Username)
		return nil
	}

	// The password is never stored, only this hash of it.
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	var superAdmin model.Role
	if err := db.Where("name = ?", model.RoleSuperAdmin).First(&superAdmin).Error; err != nil {
		return fmt.Errorf("find the super_admin role: %w", err)
	}

	user := model.AdminUser{
		Username:     admin.Username,
		Email:        adminEmail(admin.Username),
		FirstName:    "Super",
		LastName:     "Admin",
		PasswordHash: string(hash),
		Status:       model.StatusActive,
		Roles:        []model.Role{superAdmin},
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("create administrator: %w", err)
	}

	fmt.Printf("goose: created administrator %q with the super_admin role\n", admin.Username)

	return nil
}

// adminEmail is the email address of the seeded administrator. The email
// column is required and unique, and the configured username may not be an
// address, so one is made up from it to be corrected later.
func adminEmail(username string) string {
	if strings.Contains(username, "@") {
		return username
	}

	return username + "@localhost"
}
