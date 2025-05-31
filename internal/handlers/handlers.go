package handlers

import (
	"fmt"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var (
	botComandRegx = regexp.MustCompile(`^/weather\s*\[?([а-яА-ЯёЁa-zA-Z\s-]+)\]?$`)
)

const (
	comandStart = "/start"

	formatHelloMsg   = "Hello \"%s\"! These bot will sending weather of your city if you send him /weather[Moscow]"
	formatWeatherMsg = "%s: %.2f°C(%.2f°F)\n%s"
)

type WeatherResponse struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

// SendMsg is a small add-on to the regular bot.Send() to simplify things
func (s *Server) SendMsg(u *tgbotapi.Update, msg string) error {
	if _, err := s.bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg)); err != nil {
		return err
	}

	return nil
}

// StartComahd is a handler for the start command (/start), it displays a welcome
// message to the user
func (s *Server) StartComahd(u *tgbotapi.Update) error {
	if err := s.SendMsg(u, fmt.Sprintf(formatHelloMsg, u.Message.From.FirstName)); err != nil {
		s.logger.Error("err to send hello message", zap.String("whom", u.Message.From.UserName), zap.Error(err))
		return fmt.Errorf("err to send message: %w", err)
	}

	return nil
}

// GetWeather is a handler for the /weather[city] command
// it sends the user the weather in the city he sent
func (s *Server) GetWeather(u *tgbotapi.Update) error {
	city := extractCity(strings.TrimSpace(u.Message.Text))

	weather, ok, err := s.memo.Get(city, u.Message.From.LanguageCode)

	if err != nil {
		s.logger.Error("err to get weather from memory storage", zap.String("city", city), zap.Error(err))
		return err
	}

	if ok {
		if err := s.SendMsg(u, formatMsg(weather)); err != nil {
			s.logger.Error("err to send weather message", zap.String("whom", u.Message.From.UserName), zap.Error(err))

			return fmt.Errorf("err to send weather: %w", err)
		}
	}

	if !ok {
		weather, err := s.client.DoHTTP(city, u.Message.From.LanguageCode)
		if err != nil {
			s.logger.Error("err to get weather from foreighn api", zap.String("city", city), zap.Error(err))
			return err
		}

		if err := s.SendMsg(u, formatMsg(weather)); err != nil {
			s.logger.Error("err to send weather message", zap.String("whom", u.Message.From.UserName), zap.Error(err))

			return fmt.Errorf("err to send weather: %w", err)
		}

		if err := s.memo.Put(weather, u.Message.From.LanguageCode); err != nil {
			s.logger.Error("err to put weather in memory storage", zap.Any("weather", weather), zap.Error(err))
			// Don't return err
		}
	}

	return nil
}

// isForBot checks if the text of a user's message is addressed to a bot using a
// regular expression
func isForBot(msg string) bool {
	return len(botComandRegx.FindStringSubmatch(msg)) > 1
}

// extractCity takes a string, applies a regular expression to it, and returns the
// city (if any)
// Exmaple:
// input -> /weather[Moscow]
// output -> Moscow
func extractCity(msg string) string {
	return strings.TrimSpace(botComandRegx.FindStringSubmatch(msg)[1])
}

// formatMsg formats data from a Weather Response structure into a human-readable
// format using a template
func formatMsg(w WeatherResponse) string {
	return fmt.Sprintf(formatWeatherMsg, w.City, w.Temperature, toFahrenheit(w.Temperature), w.Description)
}

// toFahrenheit takes a temperature in degrees Celsius and then converts it to degrees
// Fahrenheit
func toFahrenheit(cels float64) float64 {
	return (cels * 9 / 5) + 32
}
