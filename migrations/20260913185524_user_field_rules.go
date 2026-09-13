package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"

	"xermess/internal/model"
)

func init() {
	goose.AddMigrationContext(upUserFieldRules, downUserFieldRules)
}

// upUserFieldRules gives a user field the rules its values have to keep —
// unique, a smallest and largest value, a prefix — and renames users.verified
// to email_verified, which says which address it is about.
func upUserFieldRules(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	// The rename comes first: AutoMigrate below would otherwise add the new
	// column beside the old one and leave the flags behind.
	if db.Migrator().HasColumn(&model.User{}, "verified") {
		if err := db.Migrator().RenameColumn(&model.User{}, "verified", "email_verified"); err != nil {
			return err
		}
	}

	if err := db.AutoMigrate(&model.UserField{}, &model.User{}); err != nil {
		return err
	}

	// The seeded fields say what they hold a little more plainly, and the
	// phone number shows what a rule is for.
	rules := []struct {
		name   string
		change map[string]any
	}{
		{"first_name", map[string]any{"label": "First name", "max": 100}},
		{"last_name", map[string]any{"label": "Last name", "max": 100}},
		{"phone", map[string]any{"label": "Phone number", "starts_with": "+", "max": 20, "is_unique": true}},
		{"phone_verified", map[string]any{"label": "Phone verified"}},
	}

	for _, rule := range rules {
		if err := db.Model(&model.UserField{}).
			Where("name = ?", rule.name).
			Updates(rule.change).Error; err != nil {
			return err
		}
	}

	return nil
}

// downUserFieldRules takes the rules away again and puts the flag back under
// its old name.
func downUserFieldRules(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	for _, column := range []string{"is_unique", "min", "max", "starts_with"} {
		if err := dropColumn(db, &model.UserField{}, column); err != nil {
			return err
		}
	}

	if db.Migrator().HasColumn(&model.User{}, "email_verified") {
		if err := db.Migrator().RenameColumn(&model.User{}, "email_verified", "verified"); err != nil {
			return err
		}
	}

	return nil
}

// dropColumn removes a column if it is there, so running down twice is not an
// error.
func dropColumn(db *gorm.DB, model any, column string) error {
	if !db.Migrator().HasColumn(model, column) {
		return nil
	}

	return db.Migrator().DropColumn(model, column)
}
