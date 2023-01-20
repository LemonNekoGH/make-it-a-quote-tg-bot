package sentry

import (
	"testing"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/config"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
	"github.com/getsentry/sentry-go"
	"github.com/samber/do"
	"github.com/stretchr/testify/require"
)

func TestInitSentry(t *testing.T) {
	t.Run("not initialized", func(t *testing.T) {
		require := require.New(t)
		// do empty init
		injector := do.New()
		do.Provide(injector, logger.NewLoggerService)
		// use empty sentry config
		do.Provide(injector, func(i *do.Injector) (config.ConfigService, error) {
			c := &config.Config{}
			return &config.ConfigServiceForTest{
				Conf: c,
			}, nil
		})

		impl, err := NewSentryService(injector)
		require.Empty(err)
		require.NotNil(impl)

		id := sentry.CaptureMessage("")
		require.Empty(id)
	})
}
