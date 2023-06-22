package config

import (
	"os"
	"strconv"

	"github.com/samber/do"

	_ "embed"
)

var (
	Version = "0.0.1"       // app version, will inject by -ldflags
	Env     = "development" // app environment, will inject by -ldflags
)

type TelegramConfig struct {
	Token    string `yaml:"token"`     // api token
	NotifyTo int64  `yaml:"notify-to"` // chat id will be used when notify on start
}

type SentryConfig struct {
	Dsn string `yaml:"dsn"` // sentry dsn
}

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Sentry   SentryConfig   `yaml:"sentry"`
}

type ConfigService interface {
	Config() *Config
}

type configService struct {
	conf *Config
}

type ConfigServiceForTest struct {
	Conf *Config
}

func (c *configService) Config() *Config {
	return c.conf
}

func (c *ConfigServiceForTest) Config() *Config {
	return c.Conf
}

func NewConfigService(injector *do.Injector) (ConfigService, error) {
	notifyTo, err := strconv.Atoi(os.Getenv("NOTIFY_CHAT_ID"))
	if err != nil {
		return nil, err
	}

	return &configService{
		conf: &Config{
			Telegram: TelegramConfig{
				Token:    os.Getenv("BOT_TOKEN"),
				NotifyTo: int64(notifyTo),
			},
		},
	}, nil
}
