package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/config"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/processor"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/sentry"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/telegram"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
	"github.com/samber/do"
)

func main() {
	// banner
	fmt.Printf("=================\n")
	fmt.Printf("MAKE IT A QUOTE\n")
	fmt.Printf("Version: %s\n", config.Version)
	fmt.Printf("Environment: %s\n", config.Env)
	fmt.Printf("=================\n")

	injector := do.New()
	do.Provide(injector, logger.NewLoggerService)
	do.Provide(injector, config.NewConfigService)
	do.Provide(injector, sentry.NewSentryService)
	do.MustInvoke[sentry.SentryService](injector) // 手动使依赖变得有用

	do.Provide(injector, telegram.NewTelegramBotService)
	do.Provide(injector, processor.NewProcessorsService)

	// graceful shutdown
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// start
	processors := do.MustInvoke[processor.ProcessorsService](injector)
	tg := do.MustInvoke[telegram.TelegramBotService](injector)

	tg.RegisterCommand(processors.Commands())
	go tg.Start(ctx, &wg)

	wg.Wait()
}
