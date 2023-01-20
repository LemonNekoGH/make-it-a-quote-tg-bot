package sentry

import (
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/config"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
	sentry "github.com/getsentry/sentry-go"
	"github.com/samber/do"
)

type SentryService interface{}

type sentryServiceImpl struct{}

func NewSentryService(injector *do.Injector) (SentryService, error) {
	c := do.MustInvoke[config.ConfigService](injector)
	logger := do.MustInvoke[logger.LoggerService](injector)

	if c.Config().Sentry.Dsn == "" {
		logger.Infof("没有 DSN 配置，不会初始化 Sentry")
		return nil, nil
	}
	// sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              c.Config().Sentry.Dsn,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
		Environment:      config.Env,
		Release:          config.Version,
	})
	if err != nil {
		panic(err)
	}

	logger.Infof("[成功] Sentry 初始化成功")

	return &sentryServiceImpl{}, nil
}
