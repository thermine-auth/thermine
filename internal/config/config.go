// Package config reads the settings in .env.
package config

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/spf13/viper"
)

// Config is every setting the server has.
type Config struct {
	Addr        string
	CORSOrigins []string
	DB          DB
	Admin       Admin
}

// Admin is the first administrator, created by a migration. Leaving either
// field empty means no administrator is created.
type Admin struct {
	Username string
	Password string
}

// DB is the database connection and migration settings.
type DB struct {
	Driver     string
	DSN        string
	TimeZone   string
	LogQueries bool
	Migrate    bool
	MigrateDir string
}

// read builds the reader both loaders use: .env first, then the environment,
// which wins. A missing .env is fine: in a container there are only
// environment variables.
func read() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	return v, nil
}

// LoadAdmin reads only the administrator credentials. The migration that
// seeds the first administrator uses this rather than Load, so a missing
// database setting cannot fail a migration that has no use for it.
func LoadAdmin() (Admin, error) {
	v, err := read()
	if err != nil {
		return Admin{}, err
	}

	return Admin{
		Username: v.GetString("XERMESS_ADMIN_USERNAME"),
		Password: v.GetString("XERMESS_ADMIN_PASSWORD"),
	}, nil
}

// Load reads every setting the server needs.
func Load() (Config, error) {
	v, err := read()
	if err != nil {
		return Config{}, err
	}

	v.SetDefault("XERMESS_ADDR", ":8080")
	v.SetDefault("XERMESS_CORS_ORIGINS", "http://localhost:5173")
	v.SetDefault("XERMESS_DB_DRIVER", "postgres")
	v.SetDefault("XERMESS_DB_TIMEZONE", "Asia/Bishkek")
	v.SetDefault("XERMESS_DB_LOG_QUERIES", false)
	v.SetDefault("XERMESS_DB_MIGRATE", true)
	v.SetDefault("XERMESS_DB_MIGRATE_DIR", "./migrations")

	cfg := Config{
		Addr:        v.GetString("XERMESS_ADDR"),
		CORSOrigins: splitList(v.GetString("XERMESS_CORS_ORIGINS")),
		Admin: Admin{
			Username: v.GetString("XERMESS_ADMIN_USERNAME"),
			Password: v.GetString("XERMESS_ADMIN_PASSWORD"),
		},
		DB: DB{
			Driver:     v.GetString("XERMESS_DB_DRIVER"),
			DSN:        v.GetString("XERMESS_DB_DSN"),
			TimeZone:   v.GetString("XERMESS_DB_TIMEZONE"),
			LogQueries: v.GetBool("XERMESS_DB_LOG_QUERIES"),
			Migrate:    v.GetBool("XERMESS_DB_MIGRATE"),
			MigrateDir: v.GetString("XERMESS_DB_MIGRATE_DIR"),
		},
	}

	// No default for the DSN: it names the database and carries the password,
	// so a missing one should stop the server, not quietly connect somewhere.
	if cfg.DB.DSN == "" {
		return Config{}, errors.New("XERMESS_DB_DSN is not set")
	}

	return cfg, nil
}

// splitList reads "a,b,c" into a slice.
func splitList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
