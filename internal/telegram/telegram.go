package telegram

import (
	"context"
	"fmt"
	"sync"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/config"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
)

func Init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(config.Conf.Telegram.Token)
	if err != nil {
		logger.Fatalf("TG Bot 初始化失败：%s", err.Error())
	}
	// 开发环境下开启 Debug
	bot.Debug = config.Env == "development"

	logger.Infof("[成功] TG Bot 初始化成功：%s", bot.Self.UserName)

	// send start message
	if config.Conf.Telegram.NotifyTo != 0 {
		_, err = bot.Send(tgbotapi.NewMessage(config.Conf.Telegram.NotifyTo, fmt.Sprintf("QuoteBot start success! username is @%s", bot.Self.UserName)))
		if err != nil {
			logger.Warnf("[失败] TG Bot 初始化消息发送失败：%s", err.Error())
		} else {
			logger.Infof("[成功] TG Bot 初始化消息发送成功")
		}
	}
}

// ProcessCommand 处理命令
// TODO: 完成消息处理
func processCommand(msg *tgbotapi.Message) {
	if msg == nil {
		return
	}

	toSend := tgbotapi.NewMessage(msg.Chat.ID, msg.Text)
	toSend.ReplyToMessageID = msg.MessageID

	bot.Send(toSend)
}

// WaitAndProcessCommand 等待并且处理指令
func WaitAndProcessCommand(ctx context.Context, wg *sync.WaitGroup) {
	logger.Infof("开始监听 TG 消息")
	u := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			processCommand(update.Message)
		case <-ctx.Done():
			logger.Infof("已结束消息监听")
			wg.Done()
			return
		}
	}
}
