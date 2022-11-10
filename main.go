package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/config"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/sentry"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/telegram"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
)

func main() {
	// banner
	fmt.Printf("=================\n")
	fmt.Printf("MAKE IT A QUOTE\n")
	fmt.Printf("Version: %s\n", config.Version)
	fmt.Printf("Environment: %s\n", config.Env)
	fmt.Printf("=================\n")
	// graceful shutdown
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	logger.Init()
	config.Init()
	sentry.Init()
	telegram.Init()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go telegram.WaitAndProcessCommand(ctx, &wg)

	wg.Wait()
}
