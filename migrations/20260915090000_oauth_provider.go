package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	"xermess/internal/model"
)

func init() {
	goose.AddMigrationContext(upOAuthProvider, downOAuthProvider)
}

// upOAuthProvider adds what the OAuth 2.0 / OpenID Connect endpoints need:
// the tables in model/oauth.go, a lockout for users like the administrators'
// one, and the links and registration switch an application's sign-in page
// shows.
//
// AutoMigrate adds the new columns to users and applications and leaves the
// existing ones alone. allow_registration is added by hand first: it is not
// null, so rows that already exist need a value, and the model cannot carry a
// true default without GORM refusing to store false.
func upOAuthProvider(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	statements := []string{
		`ALTER TABLE applications ADD COLUMN IF NOT EXISTS allow_registration boolean NOT NULL DEFAULT true`,
		`ALTER TABLE applications ALTER COLUMN allow_registration DROP DEFAULT`,
	}
	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			return fmt.Errorf("add allow_registration: %w", err)
		}
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Application{},
		&model.SigningKey{},
		&model.AuthorizationRequest{},
		&model.AuthorizationCode{},
		&model.RefreshToken{},
		&model.UserSession{},
		&model.PasswordReset{},
	)
	if err != nil {
		return fmt.Errorf("create provider tables: %w", err)
	}

	// A machine-to-machine application has no users to register.
	return db.Exec("UPDATE applications SET allow_registration = false WHERE type = ?", model.AppM2M).Error
}

// downOAuthProvider drops the tables and columns again. Every token issued
// stops working, since the keys that signed them are gone.
func downOAuthProvider(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}

	err = db.Migrator().DropTable(
		&model.PasswordReset{},
		&model.UserSession{},
		&model.RefreshToken{},
		&model.AuthorizationCode{},
		&model.AuthorizationRequest{},
		&model.SigningKey{},
	)
	if err != nil {
		return err
	}

	columns := []struct {
		model  any
		column string
	}{
		{&model.User{}, "LastLoginAt"},
		{&model.User{}, "FailedLoginCount"},
		{&model.User{}, "LockedUntil"},
		{&model.Application{}, "PolicyURI"},
		{&model.Application{}, "TosURI"},
		{&model.Application{}, "AllowRegistration"},
	}
	for _, c := range columns {
		if err := db.Migrator().DropColumn(c.model, c.column); err != nil {
			return err
		}
	}

	return nil
}
