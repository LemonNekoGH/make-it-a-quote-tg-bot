package processor

import (
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/telegram"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/do"
)

type ProcessorsService interface {
	Commands() telegram.CommandProcessors
}
type processorsServiceImpl struct {
	bot    telegram.TelegramBotService
	logger logger.LoggerService
}

// NewProcessorsService
func NewProcessorsService(injector *do.Injector) (ProcessorsService, error) {
	botService := do.MustInvoke[telegram.TelegramBotService](injector)
	logger := do.MustInvoke[logger.LoggerService](injector)

	p := processorsServiceImpl{
		bot:    botService,
		logger: logger,
	}

	logger.Infof("processors registered")
	return &p, nil
}

func (p *processorsServiceImpl) Commands() telegram.CommandProcessors {
	return map[string]telegram.CommandProcessor{
		"quote": p.processQuote,
	}
}

// 处理图片
// TODO:
func (p *processorsServiceImpl) processQuote(msg *tgbotapi.Message) error {
	p.logger.Infof("processing /quote command")
	_, err := p.bot.Bot().Send(tgbotapi.NewMessage(msg.Chat.ID, "currently testing"))
	return err
}
