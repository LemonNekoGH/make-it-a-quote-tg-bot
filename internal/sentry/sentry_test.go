package sentry

import (
	"testing"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/config"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
)

func TestInitSentry(t *testing.T) {
	t.Run("not initialized", func(t *testing.T) {
		require := require.New(t)
		// do empty init
		config.Conf = &config.Config{}
		Init()

		id := sentry.CaptureMessage("")
		require.Empty(id)
	})
}
