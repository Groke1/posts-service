package main

import (
	"log"
	"posts-service/internal/app"
	"posts-service/internal/config"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("can't initialize zap logger: %s", err)
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("can't initialize config: %s", err)
	}

	if err := app.Run(logger, cfg); err != nil {
		log.Fatalf("can not initialize app: %s", err)
	}
}
