package config

import (
	"io/ioutil"
	"os"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
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

var Conf *Config // 配置

// 初始化配置
func Init() {
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
	Conf = new(Config)
	err = yaml.Unmarshal(confContent, Conf)
	if err != nil {
		panic(err)
	}
	logger.Infof("配置文件加载成功")
}
