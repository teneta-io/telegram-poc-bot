package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"sync"
	"teneta-tg/internal/entities"
	"teneta-tg/internal/services"
	"teneta-tg/internal/translator"
)

type Config struct {
	Token       string
	SenderCount int
}

type Bot struct {
	mx  sync.Mutex
	ctx context.Context
	wg  *sync.WaitGroup
	api *tgbotapi.BotAPI

	messageCh      chan *MessageResponse
	userCollection map[int64]*entities.User

	userService services.UserService
	translator  *translator.Translator
}

func New(ctx context.Context, wg *sync.WaitGroup, cfg *Config, api *tgbotapi.BotAPI, userService services.UserService, translator *translator.Translator) (*Bot, error) {
	wg.Add(1)
	b := &Bot{
		ctx: ctx,
		wg:  wg,
		api: api,

		messageCh:      make(chan *MessageResponse, 200),
		userCollection: make(map[int64]*entities.User),

		userService: userService,
		translator:  translator,
	}

	for i := 0; i < cfg.SenderCount; i++ {
		go b.send()
	}

	return b, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			b.proceedUpdate(update)
		case <-b.ctx.Done():
			b.wg.Done()
			zap.S().Info("Telegram bot API stopped.")
			return
		}
	}
}

func (b *Bot) send() {
	for message := range b.messageCh {
		msg := tgbotapi.NewMessage(message.ChatId, message.Text)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

		if message.ReplyMarkup != nil {
			msg.ReplyMarkup = message.ReplyMarkup
		}

		if message.InlineMarkup != nil {
			msg.ReplyMarkup = message.InlineMarkup
		}

		msg.ParseMode = tgbotapi.ModeHTML
		msg.DisableWebPagePreview = true

		if _, err := b.api.Send(msg); err != nil {
			zap.S().Errorf("Can't send message to chat id: [%d]. Reason: %s", msg.ChatID, err.Error())
		}
	}
}

func (b *Bot) typing(ChatID int64) (tgbotapi.Message, error) {
	msg := tgbotapi.NewChatAction(ChatID, "typing")

	return b.api.Send(msg)
}

func (b *Bot) proceedUpdate(update tgbotapi.Update) {
	if update.Message != nil && update.Message.Chat != nil {
		user := b.prepareUser(update.Message)

		if update.Message.IsCommand() {
			b.proceedCommand(user, update.Message.Command())
		} else {
			b.proceedMessage(user, update.Message.Text)
		}
	}
}

func (b *Bot) response(user *entities.User, message string, args map[string]interface{}, inlineMarkup *tgbotapi.InlineKeyboardMarkup, replyMarkup *tgbotapi.ReplyKeyboardMarkup) {
	b.messageCh <- &MessageResponse{
		ChatId:       user.ChatID,
		Text:         b.translator.Translate(message, "en", args),
		InlineMarkup: inlineMarkup,
		ReplyMarkup:  replyMarkup,
	}
}

func (b *Bot) prepareUser(message *tgbotapi.Message) *entities.User {
	b.mx.Lock()
	defer b.mx.Unlock()

	if _, exist := b.userCollection[message.Chat.ID]; !exist {
		user, err := b.userService.FirstOrCreate(message.Chat.ID, message.Chat.FirstName, message.Chat.LastName, "en")
		if err != nil {
			zap.S().Errorf("can not retrieve user from db: %e", err)

			return nil
		}

		b.userCollection[user.ChatID] = user
	}

	return b.userCollection[message.Chat.ID]
}
