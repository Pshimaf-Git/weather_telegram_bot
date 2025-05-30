package handlers

import (
	"errors"
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

type TelegramHandlers struct {
	bot    *tgbotapi.BotAPI
	logger *zap.Logger
	client WeatherClient
	memo   MemoryRepo
}

func NewTGHandlers(bot *tgbotapi.BotAPI, logger *zap.Logger, client WeatherClient, memo MemoryRepo) (*TelegramHandlers, error) {
	if bot == nil || logger == nil || client == nil || memo == nil {
		return nil, errors.New("invalid input data(nil params)")
	}

	return &TelegramHandlers{
		bot:    bot,
		logger: logger,
		client: client,
		memo:   memo,
	}, nil
}

type WeatherResponse struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

type MemoryRepo interface {
	Get(city string) (WeatherResponse, bool, error)
	Put(WeatherResponse) error
}

type WeatherClient interface {
	DoHTTP(city string) (WeatherResponse, error)
}

func (t *TelegramHandlers) SendMsg(u *tgbotapi.Update, msg string) error {
	if _, err := t.bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg)); err != nil {
		return err
	}

	return nil
}

func (t *TelegramHandlers) StartComahd(u *tgbotapi.Update) error {
	if u.Message.Text != comandStart {
		t.logger.Info("comand is not start", zap.String("this comand", u.Message.Text), zap.String("waiting comand", comandStart))
		return fmt.Errorf("this comand(%s) is not start(%s)", u.Message.Text, comandStart)
	}

	if err := t.SendMsg(u, fmt.Sprintf(formatHelloMsg, u.Message.From.FirstName)); err != nil {
		t.logger.Error("err to send hello message", zap.String("whom", u.Message.From.UserName), zap.Error(err))
		return fmt.Errorf("err to send message: %w", err)
	}

	return nil
}

func (t *TelegramHandlers) GetWeather(u *tgbotapi.Update) error {
	if !isForBot(strings.TrimSpace(u.Message.Text)) {
		t.logger.Info("these message doesn't for bot", zap.String("this message", u.Message.Text))
		return fmt.Errorf("these message doesn't for bot: \"%s\"", u.Message.Text)
	}

	city := extractCity(strings.TrimSpace(u.Message.Text))

	weather, ok, err := t.memo.Get(city)

	if err != nil {
		t.logger.Error("err to get weather from memory storage", zap.String("city", city), zap.Error(err))
		return err
	}

	if ok {
		if err := t.SendMsg(u, formatMsg(weather)); err != nil {
			t.logger.Error("err to send weather message", zap.String("whom", u.Message.From.UserName), zap.Error(err))

			return fmt.Errorf("err to send weather: %w", err)
		}
	}

	if !ok {
		weather, err := t.client.DoHTTP(city)
		if err != nil {
			t.logger.Error("err to get weather from foreighn api", zap.String("city", city), zap.Error(err))
			return err
		}

		if err := t.SendMsg(u, formatMsg(weather)); err != nil {
			t.logger.Error("err to send weather message", zap.String("whom", u.Message.From.UserName), zap.Error(err))

			return fmt.Errorf("err to send weather: %w", err)
		}

		if err := t.memo.Put(weather); err != nil {
			t.logger.Error("err to put weather in memory storage", zap.Any("weather", weather), zap.Error(err))
			// Don't return err
		}
	}

	return nil
}

func isForBot(msg string) bool {
	return len(botComandRegx.FindStringSubmatch(msg)) > 1
}

func extractCity(msg string) string {
	return strings.TrimSpace(botComandRegx.FindStringSubmatch(msg)[1])
}

func formatMsg(w WeatherResponse) string {
	return fmt.Sprintf(formatWeatherMsg, w.City, w.Temperature, toFahrenheit(w.Temperature), w.Description)
}

func toFahrenheit(cels float64) float64 {
	return (cels * 9 / 5) + 32
}
