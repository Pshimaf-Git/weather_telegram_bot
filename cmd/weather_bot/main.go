package main

import (
	"fmt"
	"log"
	"weather_bot/config"
	"weather_bot/internal/handlers"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	cfg, err := config.Load("CONFIG_PATH")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	bot, err := cfg.NewBot()
	if err != nil {
		return fmt.Errorf("create telegram bot: %w", err)
	}

	client, err := cfg.NewWeatherClient()
	if err != nil {
		return fmt.Errorf("create new weather client: %w", err)
	}

	repo, err := cfg.NewRepo()
	if err != nil {
		return fmt.Errorf("err init connect to repository: %w", err)
	}

	logger, err := cfg.NewLogger()
	if err != nil {
		return fmt.Errorf("create zap logger: %w", err)
	}

	server, err := handlers.NewServer(bot, logger, client, repo)
	if err != nil {
		return fmt.Errorf("init server: %w", err)
	}

	return fmt.Errorf("server running error: %w", server.Run())
}
