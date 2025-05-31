package main

import (
	"log"
	"weather_bot/config"
	"weather_bot/internal/handlers"
)

func main() {
	cfg, err := config.Load("CONFIG_PATH")
	if err != nil {
		log.Fatalf("err to load config: %s", err.Error())
	}

	bot, err := cfg.NewBot()
	if err != nil {
		log.Fatalf("err to create telegram bot: %s", err.Error())
	}

	client, err := cfg.NewWeatherClient()
	if err != nil {
		log.Fatalf("err to create new weather client: %s", err.Error())
	}

	repo, err := cfg.NewRepo()
	if err != nil {
		log.Fatalf("err connect  to repository: %s", err.Error())
	}

	logger, err := cfg.NewLogger()
	if err != nil {
		log.Fatalf("err to create zap logger: %s", err.Error())
	}

	server, err := handlers.NewServer(bot, logger, client, repo)
	if err != nil {
		log.Fatalf("err to init server: %s", err.Error())
	}

	if err := server.Run(); err != nil {
		log.Fatalf("err server running: %s", err.Error())
	}
}
