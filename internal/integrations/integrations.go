// integrations шы an intermediate layer where it is determined which
// implementation of the WeatherClient interface we will use
package integrations

import (
	"errors"
	"fmt"
	"strings"

	"weather_bot/internal/handlers"
	weather "weather_bot/internal/integrations/open_weather"
)

const (
	openWeather = "open_weather"
)

var (
	ErrUnkown = errors.New("unknown integration")
)

// New return a new instance of the WeatherClient interface
func New(name string, apiKey string, baseURl string, defaultLang string) (handlers.WeatherClient, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case openWeather:
		return weather.New(apiKey, baseURl, defaultLang)
	default:
	}

	return nil, fmt.Errorf("%w: %s", ErrUnkown, name)
}
