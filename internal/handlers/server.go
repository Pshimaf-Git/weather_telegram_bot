package handlers

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var (
	ErrServerStoped = errors.New("server stoped")
)

type MemoryRepo interface {
	Get(city string, lang string) (WeatherResponse, bool, error)
	Put(weather WeatherResponse, lang string) error
}

type WeatherClient interface {
	DoHTTP(city string, lang string) (WeatherResponse, error)
}

type Server struct {
	bot    *tgbotapi.BotAPI
	logger *zap.Logger
	client WeatherClient
	memo   MemoryRepo
}

func NewServer(bot *tgbotapi.BotAPI, logger *zap.Logger, client WeatherClient, memo MemoryRepo) (*Server, error) {
	if bot == nil || logger == nil || client == nil || memo == nil {
		return nil, errors.New("invalid input data(nil params)")
	}

	return &Server{
		bot:    bot,
		logger: logger,
		client: client,
		memo:   memo,
	}, nil
}

func (s *Server) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := s.bot.GetUpdatesChan(u)
	for update := range updates {
		switch {
		case update.Message.Text == comandStart:
			s.StartComahd(&update)
		case isForBot(update.Message.Text):
			s.GetWeather(&update)
		}
	}

	return ErrServerStoped
}
