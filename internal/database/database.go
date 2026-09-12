// Package database opens the database connection.
package database

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"xermess/internal/config"
)

// Open connects to Postgres and checks the connection works, so a bad
// database stops the server at startup instead of on the first request.
func Open(cfg config.DB) (*gorm.DB, error) {
	if cfg.Driver != "postgres" {
		return nil, fmt.Errorf("unsupported database driver %q", cfg.Driver)
	}

	level := logger.Warn
	if cfg.LogQueries {
		level = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn(cfg)), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

// Close shuts the connection pool down.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// dsn adds the time zone to the DSN when it does not already carry one. It is
// added by hand because escaping would turn the slash in "Asia/Bishkek" into
// %2F, which Postgres rejects.
func dsn(cfg config.DB) string {
	if strings.Contains(cfg.DSN, "TimeZone=") {
		return cfg.DSN
	}

	separator := "?"
	if strings.Contains(cfg.DSN, "?") {
		separator = "&"
	}

	return cfg.DSN + separator + "TimeZone=" + cfg.TimeZone
}
