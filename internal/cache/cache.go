package cache

import (
	"fmt"
	"weather_bot/internal/handlers"
)

type Cahche struct{}

func (c *Cahche) Get(city string, lang string) (handlers.WeatherResponse, bool, error) {
	return handlers.WeatherResponse{}, false, nil
}

func (c *Cahche) Put(weather handlers.WeatherResponse, lang string) error {
	fmt.Printf("Save to cache: %v\n", weather)
	return nil
}

func New(port string) (handlers.MemoryRepo, error) {
	return &Cahche{}, nil
}
