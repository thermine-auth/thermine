// Command xermess runs the authentication server.
package main

import (
	"log/slog"
	"os"

	"xermess/internal/config"
	"xermess/internal/database"
	"xermess/internal/server"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := run(log); err != nil {
		log.Error("fatal", "error", err)
		os.Exit(1)
	}
}

// run is the whole startup, in order: read the configuration, open the
// database, apply migrations, serve. It is separate from main so every step
// can return an error instead of exiting from the middle of the startup.
func run(log *slog.Logger) error {
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

	// Run blocks until the server stops. Gin listens for us: Run is a wrapper
	// around net/http's ListenAndServe.
	return server.New(cfg, log).Run(cfg.Addr)
}
