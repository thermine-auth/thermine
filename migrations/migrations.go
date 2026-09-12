// Package migrations holds the database migrations, one file per migration,
// each registering itself with goose from its init function.
//
// Importing this package is what makes the migrations exist; internal/database
// imports it for that reason alone.
package migrations

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// gormTx wraps the transaction goose hands a migration in GORM, so a
// migration can use the model structs and the migrator instead of writing DDL
// by hand. Everything it does runs in goose's transaction, so a migration
// that fails half way leaves nothing behind.
func gormTx(tx *sql.Tx) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: tx}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("wrap transaction: %w", err)
	}
	return db, nil
}
