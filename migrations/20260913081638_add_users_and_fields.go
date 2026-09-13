package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	"xermess/internal/model"
)

func init() {
	goose.AddMigrationContext(upAddUsersAndFields, downAddUsersAndFields)
}

// upAddUsersAndFields creates the users table and the table describing its
// user-defined fields, then seeds the fields an organisation is most likely
// to want. The seeded ones are ordinary rows: they can be removed from the
// panel like any other.
func upAddUsersAndFields(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.UserField{}, &model.User{}); err != nil {
		return err
	}

	fields := []model.UserField{
		{Name: "first_name", Label: "First name", Type: model.FieldText, Position: 1},
		{Name: "last_name", Label: "Last name", Type: model.FieldText, Position: 2},
		{Name: "phone", Label: "Phone", Type: model.FieldText, Position: 3},
		{Name: "phone_verified", Label: "Phone verified", Type: model.FieldBool, Position: 4},
	}

	for _, field := range fields {
		var existing int64
		if err := db.Model(&model.UserField{}).Where("name = ?", field.Name).Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			continue
		}

		if err := db.Create(&field).Error; err != nil {
			return fmt.Errorf("seed field %s: %w", field.Name, err)
		}
	}

	return nil
}

// downAddUsersAndFields drops both tables, and with them every user.
func downAddUsersAndFields(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	return db.Migrator().DropTable(&model.User{}, &model.UserField{})
}
