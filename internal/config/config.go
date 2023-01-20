package config

import (
	"io/ioutil"
	"os"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
	"github.com/samber/do"
	"gopkg.in/yaml.v3"

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
	logger := do.MustInvoke[logger.LoggerService](injector)
	var confContent []byte
	var err error = nil
	// 读取环境变量
	configFile := os.Getenv("BOT_CONFIG_PATH")

	var file *os.File
	// 读取环境变量中指定的文件
	file, err = os.Open(configFile)
	if err != nil {
		panic(err)
	}
	confContent, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// 读取完成，映射成对象
	conf := new(Config)
	err = yaml.Unmarshal(confContent, conf)
	if err != nil {
		panic(err)
	}
	logger.Infof("配置文件加载成功")

	return &configService{
		conf: conf,
	}, nil
}
