// Command xermess runs the authentication server.
package main

import (
	"log/slog"
	"os"

	"xermess/internal/api"
	"xermess/internal/config"
	"xermess/internal/database"
	"xermess/internal/store"
)

// version and commit are stamped in at build time by scripts/build.sh and
// the Dockerfile, so a running server can say which build it is. A plain
// `go run` leaves them as they are.
var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := run(log); err != nil {
		log.Error("fatal", "error", err)
		os.Exit(1)
	}
}

// run is the whole startup, in order: read the configuration, open the
// database, apply migrations, then serve on top of a store. It is separate
// from main so every step can return an error instead of exiting from the
// middle of the startup.
func run(log *slog.Logger) error {
	log.Info("xermess starting", "version", version, "commit", commit)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := database.Open(cfg.DB)
	if err != nil {
		return err
	}
	defer database.Close(db)

	log.Info("database connected")

	if cfg.DB.Migrate {
		if err := database.Migrate(db, cfg.DB, log); err != nil {
			return err
		}
	}

	log.Info("server listening", "addr", cfg.Addr)

	// The store is the only thing that queries the database; the server is
	// handed that rather than the connection itself.
	//
	// Run blocks until the server stops. Gin listens for us: Run is a wrapper
	// around net/http's ListenAndServe.
	return api.New(cfg, store.New(db), log).Run(cfg.Addr)
}
