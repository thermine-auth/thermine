package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// clearEnv unsets every XERMESS_ variable for the duration of the test, so a
// developer's own environment cannot change the result. Viper treats an empty
// variable as unset.
func clearEnv(t *testing.T) {
	t.Helper()

	for _, entry := range os.Environ() {
		if name, _, found := strings.Cut(entry, "="); found && strings.HasPrefix(name, "XERMESS_") {
			t.Setenv(name, "")
		}
	}
}

// writeEnv puts a .env in a temporary directory and makes it the working
// directory, because Load reads ./.env.
func writeEnv(t *testing.T, contents string) {
	t.Helper()

	clearEnv(t)

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
}

func TestLoadReadsEnvFile(t *testing.T) {
	writeEnv(t, `
XERMESS_ADDR=:9000
XERMESS_CORS_ORIGINS=http://a.test, http://b.test
XERMESS_DB_DSN=postgres://user:pw@localhost:5432/mydb
XERMESS_DB_LOG_QUERIES=true
XERMESS_DB_MIGRATE=false
`)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Addr != ":9000" {
		t.Errorf("Addr = %q, want :9000", cfg.Addr)
	}
	if cfg.DB.DSN != "postgres://user:pw@localhost:5432/mydb" {
		t.Errorf("DSN = %q", cfg.DB.DSN)
	}
	if !cfg.DB.LogQueries {
		t.Error("LogQueries = false, want true")
	}
	if cfg.DB.Migrate {
		t.Error("Migrate = true, want false")
	}

	want := []string{"http://a.test", "http://b.test"}
	if len(cfg.CORSOrigins) != len(want) {
		t.Fatalf("CORSOrigins = %v, want %v", cfg.CORSOrigins, want)
	}
	for i, origin := range want {
		if cfg.CORSOrigins[i] != origin {
			t.Errorf("CORSOrigins[%d] = %q, want %q (spaces should be trimmed)", i, cfg.CORSOrigins[i], origin)
		}
	}
}

func TestLoadUsesDefaults(t *testing.T) {
	writeEnv(t, "XERMESS_DB_DSN=postgres://u:p@localhost:5432/db\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Addr != ":8080" {
		t.Errorf("Addr = %q, want the default :8080", cfg.Addr)
	}
	if cfg.DB.Driver != "postgres" {
		t.Errorf("Driver = %q, want the default postgres", cfg.DB.Driver)
	}
	if cfg.DB.TimeZone != "Asia/Bishkek" {
		t.Errorf("TimeZone = %q, want the default Asia/Bishkek", cfg.DB.TimeZone)
	}
	if !cfg.DB.Migrate {
		t.Error("Migrate = false, want the default true")
	}
	if cfg.DB.MigrateDir != "./migrations" {
		t.Errorf("MigrateDir = %q, want the default ./migrations", cfg.DB.MigrateDir)
	}
}

// The environment has to win over .env, or a container could not change a
// single setting without rewriting the file.
func TestEnvironmentOverridesEnvFile(t *testing.T) {
	writeEnv(t, "XERMESS_ADDR=:8080\nXERMESS_DB_DSN=postgres://u:p@localhost:5432/from_file\n")
	t.Setenv("XERMESS_ADDR", ":7777")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Addr != ":7777" {
		t.Errorf("Addr = %q, want :7777 from the environment", cfg.Addr)
	}
}

func TestLoadFailsWithoutDSN(t *testing.T) {
	writeEnv(t, "XERMESS_ADDR=:8080\n")

	if _, err := Load(); err == nil {
		t.Fatal("want an error when XERMESS_DB_DSN is missing, got none")
	}
}

// A missing .env is normal in a container, where only real environment
// variables are set.
func TestLoadWithoutEnvFile(t *testing.T) {
	clearEnv(t)
	t.Chdir(t.TempDir())
	t.Setenv("XERMESS_DB_DSN", "postgres://u:p@localhost:5432/db")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Addr != ":8080" {
		t.Errorf("Addr = %q, want the default :8080", cfg.Addr)
	}
}
