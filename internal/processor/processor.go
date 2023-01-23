package processor

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"golang.org/x/image/draw"

	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/internal/telegram"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/logger"
	"github.com/LemonNekoGH/make-it-a-quote-tg-bot/pkg/texttoimage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/samber/do"

	_ "embed"
)

type ProcessorsService interface {
	Commands() telegram.CommandProcessors
}
type processorsServiceImpl struct {
	bot           telegram.TelegramBotService
	logger        logger.LoggerService
	font          *truetype.Font
	mask          image.Image
	defaultAvatar image.Image
}

//go:embed AlibabaPuHui.ttf
var fontAlibaba []byte

//go:embed default-avatar.png
var defaultAvatar []byte

//go:embed gradient-mask.png
var maskFile []byte

// NewProcessorsService
func NewProcessorsService(injector *do.Injector) (ProcessorsService, error) {
	botService := do.MustInvoke[telegram.TelegramBotService](injector)
	logger := do.MustInvoke[logger.LoggerService](injector)

	// load font
	font, err := freetype.ParseFont(fontAlibaba)
	if err != nil {
		logger.Fatalf("processors register failed", err.Error())
	}
	// load avatar
	avatar, err := png.Decode(bytes.NewBuffer(defaultAvatar))
	if err != nil {
		logger.Fatalf("processors register failed", err.Error())
	}
	// load mask
	mask, err := png.Decode(bytes.NewBuffer(maskFile))
	if err != nil {
		logger.Fatalf("processors register failed", err.Error())
	}

	logger.Infof("processors registered")
	return &processorsServiceImpl{
		bot:           botService,
		logger:        logger,
		font:          font,
		defaultAvatar: avatar,
		mask:          mask,
	}, nil
}

func (p *processorsServiceImpl) Commands() telegram.CommandProcessors {
	return map[string]telegram.CommandProcessor{
		"quote": p.processQuote,
	}
}

// 处理图片
func (p *processorsServiceImpl) processQuote(msg *tgbotapi.Message) error {
	p.logger.Infof("processing /quote command")

	replyedMsg := msg.ReplyToMessage
	if replyedMsg == nil {
		// no reply
		p.logger.Infof("no replyed message, id: %d, chat id: %d", msg.MessageID, msg.Chat.ID)
		replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "You need to reply a message")
		replyMsg.ReplyToMessageID = msg.MessageID
		_, err := p.bot.Bot().Send(replyMsg)
		p.logger.Errorf("send error message failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
		return err
	}
	if replyedMsg.Text == "" {
		// reply no content
		p.logger.Infof("message no content, id: %d, chat id: %d", msg.MessageID, msg.Chat.ID)
		replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "This message has no text content")
		replyMsg.ReplyToMessageID = msg.MessageID
		_, err := p.bot.Bot().Send(replyMsg)
		if err != nil {
			p.logger.Errorf("send error message failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
		}
		return err
	}
	// get username and avatar
	if replyedMsg.From == nil {
		p.logger.Infof("sender is nil, id: %d, chat id: %d", msg.MessageID, msg.Chat.ID)
		replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Message sender is nil, it might be from a channel")
		replyMsg.ReplyToMessageID = msg.MessageID
		_, err := p.bot.Bot().Send(replyMsg)
		p.logger.Errorf("send error message failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
		return err
	}
	senderName := ""
	userId := int64(0)
	if replyedMsg.ForwardFrom != nil {
		// forwarded message
		senderName = replyedMsg.ForwardFrom.UserName
		userId = replyedMsg.ForwardFrom.ID
	} else if replyedMsg.ForwardSenderName != "" {
		// forwarded message, but privacy limited
		p.logger.Infof("forward message privacy limited, id: %d, chat id: %d", msg.MessageID, msg.Chat.ID)
		replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Message was forwarded from a user, but cannot get username because user's privacy limit")
		replyMsg.ReplyToMessageID = msg.MessageID
		_, err := p.bot.Bot().Send(replyMsg)
		return err
	} else {
		senderName = replyedMsg.From.UserName
		if senderName == "" {
			// empty username
			senderName = "no_username"
		}
		userId = replyedMsg.From.ID
	}

	// send processing message
	replyMsg := tgbotapi.NewMessage(msg.Chat.ID, "Processing...")
	replyMsg.ReplyToMessageID = msg.MessageID
	processingMsg, err := p.bot.Bot().Send(replyMsg)
	defer func() {
		// edit msg
		if err != nil {
			params := tgbotapi.Params{}
			params.AddNonZero64("chat_id", replyMsg.ChatID)
			params.AddNonZero("message_id", processingMsg.MessageID)
			params["text"] = "Process failed"
			resp, err2 := p.bot.Bot().MakeRequest("editMessageText", params)
			if err2 != nil {
				p.logger.Errorf("edit process message failed: %s", err2.Error())
			}
			if resp != nil || !resp.Ok {
				p.logger.Errorf("edit process message failed: %s", resp.Description)
			}
		}
	}()
	if err != nil {
		p.logger.Errorf("send processing msg failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
		return err
	}
	p.logger.Infof("processing msg sent, id: %d", processingMsg.MessageID)

	// get avatar
	photos, err := p.bot.Bot().GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: userId,
	})
	if err != nil {
		p.logger.Errorf("get user avatar failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
		return err
	}
	avatar := []byte{}
	if len(photos.Photos) != 0 {
		// user has photo, use index 0
		photoId := photos.Photos[0][0].FileID
		fileUrl, err := p.bot.Bot().GetFileDirectURL(photoId)
		if err != nil {
			p.logger.Errorf("get avatar url failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
			return err
		}
		// download file
		resp, err := http.Get(fileUrl)
		if err != nil {
			p.logger.Errorf("download file failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
			return err
		}
		defer resp.Body.Close()

		buf := bytes.NewBuffer([]byte{})
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			p.logger.Errorf("copy avatar response failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
			return err
		}
		avatar = buf.Bytes()
	}

	// generate text image
	options := &texttoimage.TextToImageOptions{
		Font:     p.font,
		FontSize: 12,
		DPI:      144,
		Padding:  20,
		MaxWidth: 600,
	}
	textImg := texttoimage.TextToImage("「"+replyedMsg.Text+"」", options)
	// generate username image
	usernameImg := texttoimage.TextToImage("@"+senderName, options)
	// combine text and username TODO: can be optimized to reduce a image
	width, height := textImg.Rect.Dx(), textImg.Rect.Dy()+usernameImg.Rect.Dy()
	p.logger.Debugf("text size: %dx%d, username size: %dx%d", textImg.Rect.Dx(), textImg.Rect.Dy(), usernameImg.Rect.Dx(), usernameImg.Rect.Dy())
	p.logger.Debugf("text and username size: %dx%d", width, height)
	textAndUsername := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(textAndUsername, textImg.Bounds(), textImg, image.Point{}, draw.Over)
	draw.Draw(textAndUsername, image.Rect(0, textImg.Rect.Dy(), width, height), usernameImg, image.Point{}, draw.Over)
	// sacle avatar
	avatarImg := p.defaultAvatar
	if len(avatar) != 0 {
		avatarImg, err = jpeg.Decode(bytes.NewBuffer(avatar)) // TODO: can load in initial func
		if err != nil {
			p.logger.Errorf("decode avatar failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
			return err
		}
	}
	avatarAfterScale := image.NewRGBA(image.Rect(0, 0, height, height))
	draw.ApproxBiLinear.Scale(avatarAfterScale, avatarAfterScale.Rect, avatarImg, avatarImg.Bounds(), draw.Over, nil)
	// combine text, username and avatar
	result := image.NewRGBA(image.Rect(0, 0, avatarAfterScale.Rect.Dx()+textAndUsername.Rect.Dx(), height))
	draw.Draw(result, avatarAfterScale.Bounds(), avatarAfterScale, image.Point{}, draw.Over)
	draw.Draw(result, image.Rect(height, 0, width+height, height), textAndUsername, image.Point{}, draw.Over)
	// add mask
	mask, err := png.Decode(bytes.NewBuffer(maskFile)) // TODO: can load in initial func
	if err != nil {
		p.logger.Errorf("decode mask failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
		return err
	}
	maskAfterScale := image.NewRGBA(image.Rect(0, 0, height, height))
	draw.ApproxBiLinear.Scale(maskAfterScale, maskAfterScale.Rect, mask, mask.Bounds(), draw.Over, nil)
	draw.Draw(result, maskAfterScale.Bounds(), maskAfterScale, image.Point{}, draw.Over)
	// upload file
	resultBuff := bytes.NewBuffer([]byte{})
	err = png.Encode(resultBuff, result)
	if err != nil {
		p.logger.Errorf("encode image failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
		return err
	}
	p.logger.Infof("Image generated, id: %d, chat id: %d", msg.MessageID, msg.Chat.ID)

	photoCfg := tgbotapi.NewPhoto(msg.Chat.ID, tgbotapi.FileBytes{
		Name:  "generated",
		Bytes: resultBuff.Bytes(),
	})
	photoCfg.ReplyToMessageID = msg.MessageID
	_, err = p.bot.Bot().Send(photoCfg)
	if err != nil {
		p.logger.Errorf("send image failed, err: %s, id: %d, chat id: %d", err.Error(), msg.MessageID, msg.Chat.ID)
		return err
	}
	p.logger.Infof("Image sent, id: %d, chat id: %d", msg.MessageID, msg.Chat.ID)
	// delete message
	params := tgbotapi.Params{}
	params.AddNonZero64("chat_id", replyMsg.ChatID)
	params.AddNonZero("message_id", processingMsg.MessageID)
	resp, err2 := p.bot.Bot().MakeRequest("deleteMessage", params)
	if err2 != nil {
		p.logger.Errorf("delete message failed: %s", err2.Error())
	}
	if resp != nil || !resp.Ok {
		p.logger.Errorf("delete message failed: %d", resp.ErrorCode)
	}
	p.logger.Infof("message deleted, id: %d, chat id: %d", msg.MessageID, msg.Chat.ID)
	return err
}
