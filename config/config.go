package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"weather_bot/internal/cache"
	"weather_bot/internal/handlers"
	"weather_bot/internal/integrations"
)

const (
	dev         = "dev"
	prod        = "prod"
	production  = "production"
	development = "development"
)

type ServerCfg struct {
	BotToken      string `yaml:"bot_token"`
	Logger        `yaml:"logger"`
	WeatherClient `yaml:"weather_client"`
	Repository    `yaml:"repository"`
}

type WeatherClient struct {
	Name        string `yaml:"client_name"`
	BaseURL     string `yaml:"base_url"`
	ApiKey      string `yaml:"api_key"`
	DefaultLang string `yaml:"default_lang"`
}

type Logger struct {
	LoggerName string `yaml:"logger_name"`
}

type Repository struct {
	Port string `yaml:"port"`
}

// Load takes data from the file whose path is specified in the environment variabl
// (its name is the input parameter)
func Load(cfgEnv string) (*ServerCfg, error) {
	cfgPath := os.Getenv(cfgEnv)
	if strings.EqualFold(cfgPath, "") {
		return nil, fmt.Errorf("err: empty config path")
	}

	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("err to open config file<%s>: %w", cfgPath, err)
	}

	defer cfgFile.Close()

	cfgData, err := io.ReadAll(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("err to read config file<%s>: %w", cfgPath, err)
	}

	var cfg ServerCfg
	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		return nil, fmt.Errorf("err to unmarshal data in config file<%s>: %W", cfgPath, err)
	}

	return &cfg, nil
}

// NewBot takes bot token from config and create bot with it
func (cfg *ServerCfg) NewBot() (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPI(cfg.BotToken)
}

// NewWeatherClient
func (cfg *ServerCfg) NewWeatherClient() (handlers.WeatherClient, error) {
	return integrations.New(cfg.Name, cfg.ApiKey, cfg.BaseURL, cfg.DefaultLang)
}

func (cfg *ServerCfg) NewRepo() (handlers.MemoryRepo, error) {
	return cache.New(cfg.Port)
}

func (cfg *ServerCfg) NewLogger(options ...zap.Option) (*zap.Logger, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.LoggerName)) {
	case development, dev:
		return zap.NewDevelopment(options...)
	case production, prod:
		return zap.NewProduction(options...)
	default:
	}

	return nil, fmt.Errorf("unknown logger name: %s", cfg.LoggerName)
}
