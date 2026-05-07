package main

import (
	"log"

	"github.com/71g3pf4c3/gophermart/config"
	"github.com/71g3pf4c3/gophermart/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
