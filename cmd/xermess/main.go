package main

import (
	"log"

	"xermess/internal/config"
	"xermess/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := server.New(cfg).Run(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}
