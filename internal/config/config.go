package config

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Addr        string
	CORSOrigins []string
}

func Load() (Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	v.SetDefault("XERMESS_ADDR", ":8080")
	v.SetDefault("XERMESS_CORS_ORIGINS", "http://localhost:5173")

	if err := v.ReadInConfig(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Config{}, fmt.Errorf("read .env: %w", err)
	}

	return Config{
		Addr:        v.GetString("XERMESS_ADDR"),
		CORSOrigins: splitList(v.GetString("XERMESS_CORS_ORIGINS")),
	}, nil
}

func splitList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
