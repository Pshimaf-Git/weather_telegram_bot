// weather is a package that provides one
// implementation of the WeatherClient interface for
// queries to the open weather service
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

const (
	UnknownCity = "unkown city"
)

type OpenWeatheClient struct {
	baseURL     string
	apiKey      string
	defaultLang string
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

// New returns a new instance of the OpenWeatheClient structure
func New(apiKey string, baseURL string, defaultLang string) (*OpenWeatheClient, error) {
	if strings.EqualFold(baseURL, "") {
		return nil, errors.New("emprty base url")
	}

	return &OpenWeatheClient{
		apiKey:      apiKey,
		baseURL:     baseURL,
		defaultLang: defaultLang,
	}, nil
}

// DoHTTP makes an http request to get the weather in the given city
// returning the result in the given language
func (c *OpenWeatheClient) DoHTTP(city string, lang string) (handlers.WeatherResponse, error) {
	if strings.EqualFold(lang, "") {
		lang = c.defaultLang
	}

	fullURL := c.makeHTTP(city, lang)

	resp, err := http.Get(fullURL)
	if err != nil {
		return handlers.WeatherResponse{}, fmt.Errorf("err to get response(url=%s): %w", fullURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return handlers.WeatherResponse{
				City:        city,
				Temperature: 0.0,
				Description: UnknownCity,
			}, nil
		}

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

// makeHTTP is a helper method that builds an http request with a given base path and
// additional parameters
func (c *OpenWeatheClient) makeHTTP(city string, lang string) string {
	params := url.Values{}
	params.Add("q", city)
	params.Add("appid", c.apiKey)
	params.Add("units", "metric")
	params.Add("lang", lang)

	return c.baseURL + "?" + params.Encode()
}
