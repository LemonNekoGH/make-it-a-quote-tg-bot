package telegram

import (
	"context"
	"fmt"
	"sync"

	_ "embed"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/config"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/do"
	"golang.org/x/exp/maps"
)

type TelegramBotService interface {
	// 启动
	Start(ctx context.Context, wg *sync.WaitGroup)
	// 注册命令处理器
	RegisterCommand(processors CommandProcessors)
	// 获取机器人实例
	Bot() *tgbotapi.BotAPI
}

type CommandProcessor func(msg *tgbotapi.Message) error

// 命令处理器组
type CommandProcessors map[string]CommandProcessor

type telgeramBotServiceImpl struct {
	bot               *tgbotapi.BotAPI
	commandProcessors map[string]CommandProcessor
	logger            logger.LoggerService
	config            config.ConfigService
}

// Start 请使用 go 语法运行此协程
func (t *telgeramBotServiceImpl) Start(ctx context.Context, wg *sync.WaitGroup) {
	// send start message
	if t.config.Config().Telegram.NotifyTo != 0 {
		_, err := t.bot.Send(tgbotapi.NewMessage(t.config.Config().Telegram.NotifyTo, fmt.Sprintf("QuoteBot started! username is @%s", t.bot.Self.UserName)))
		if err != nil {
			t.logger.Warnf("[失败] TG Bot 初始化消息发送失败：%s", err.Error())
		} else {
			t.logger.Infof("[成功] TG Bot 初始化消息发送成功")
		}
	}
	defer wg.Done()

	t.logger.Infof("Start to receive telegram message")
	u := tgbotapi.NewUpdate(0)
	updates := t.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			t.processCommand(update.Message)
		case <-ctx.Done():
			t.logger.Infof("Stopped")
			return
		}
	}
}

func (t *telgeramBotServiceImpl) RegisterCommand(processors CommandProcessors) {
	maps.Copy(t.commandProcessors, processors)
}

func (t *telgeramBotServiceImpl) Bot() *tgbotapi.BotAPI {
	return t.bot
}

// ProcessCommand 处理命令
func (t *telgeramBotServiceImpl) processCommand(msg *tgbotapi.Message) {
	if msg == nil {
		return
	}

	t.logger.Debugf("received message: %s, is command: %v, command: %s", msg.Text, msg.IsCommand(), msg.Command())
	// 不是命令不处理
	if !msg.IsCommand() {
		return
	}
	// 不是发送给自己的不处理
	if !t.bot.IsMessageToMe(*msg) {
		return
	}

	// 把命令路由到该去的地方
	cmd := msg.Command()

	processor := t.commandProcessors[cmd]
	if processor == nil {
		t.logger.Debugf("no such command: %s", cmd)
	} else {
		err := processor(msg)
		if err != nil {
			t.logger.Errorf("processer error: %s, command: %s", err.Error(), cmd)
		}

		// recover panic
		anyErr := recover()
		if anyErr != nil {
			t.logger.Errorf("processer error: %s, command: %s", anyErr, cmd)
		}
	}
}

func NewTelegramBotService(injector *do.Injector) (TelegramBotService, error) {
	c := do.MustInvoke[config.ConfigService](injector)
	logger := do.MustInvoke[logger.LoggerService](injector)

	bot, err := tgbotapi.NewBotAPI(c.Config().Telegram.Token)
	if err != nil {
		logger.Fatalf("TG Bot 初始化失败：%s", err.Error())
	}

	logger.Infof("[成功] TG Bot 初始化成功：%s", bot.Self.UserName)

	return &telgeramBotServiceImpl{
		bot:               bot,
		config:            c,
		commandProcessors: map[string]CommandProcessor{},
		logger:            logger,
	}, nil
}
