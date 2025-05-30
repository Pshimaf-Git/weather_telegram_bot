package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"weather_bot/internal/handlers"
)

type OpenWeatheClient struct {
	baseURL string
	apiKey  string
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	}
}

func New(apiKey string, baseURL string) (*OpenWeatheClient, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("emprty api")
	}

	return &OpenWeatheClient{
		apiKey:  apiKey,
		baseURL: baseURL,
	}, nil
}

func (c *OpenWeatheClient) DoHTTP(city string) (handlers.WeatherResponse, error) {
	fullURL := makeHTTP(city, c.apiKey, c.baseURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return handlers.WeatherResponse{}, fmt.Errorf("err to get response(url=%s): %w", fullURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handlers.WeatherResponse{}, fmt.Errorf("error HTTP status code: %d", resp.StatusCode)
	}

	var weather weatherData
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return handlers.WeatherResponse{}, fmt.Errorf("err to decoded response body(url=%s): %w", fullURL, err)
	}

	desc := ""
	if len(weather.Weather) > 0 {
		desc = weather.Weather[0].Description
	}

	return handlers.WeatherResponse{
		City:        weather.Name,
		Temperature: weather.Main.Temp,
		Description: desc,
	}, nil
}

func makeHTTP(city string, apiKey string, baseURL string) string {
	params := url.Values{}
	params.Add("q", city)
	params.Add("appid", apiKey)
	params.Add("units", "metric")
	params.Add("lang", "en")

	return baseURL + "?" + params.Encode()
}
