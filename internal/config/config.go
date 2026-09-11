// Package config loads runtime configuration from the environment.
package config

import (
	"os"
	"time"
)

// Config holds the settings the auth server needs to run.
type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}

// Load reads configuration from the environment, falling back to defaults.
func Load() Config {
	return Config{
		Addr:            env("XERMESS_ADDR", ":8080"),
		ShutdownTimeout: 10 * time.Second,
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
