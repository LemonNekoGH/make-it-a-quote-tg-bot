package sentry

import (
	"log"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/config"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
	sentry "github.com/getsentry/sentry-go"
)

func Init() {
	if config.Conf.Sentry.Dsn == "" {
		log.Println("没有 DSN 配置，不会初始化 Sentry")
		return
	}
	// sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.Conf.Sentry.Dsn,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
		Environment:      config.Env,
		Release:          config.Version,
	})
	if err != nil {
		panic(err)
	}

	logger.Infof("[成功] Sentry 初始化成功")
}
