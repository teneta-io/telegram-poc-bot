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

	msgCh          chan *tgbotapi.MessageConfig
	userCollection map[int64]*entities.User

	commandManager *CommandManager
	messageManager *MessageManager

	userService services.UserService
	translator  *translator.Translator
}

func New(ctx context.Context, wg *sync.WaitGroup, cfg *Config, api *tgbotapi.BotAPI, commandManager *CommandManager, messageManager *MessageManager, userService services.UserService, translator *translator.Translator) (*Bot, error) {
	wg.Add(1)
	b := &Bot{
		ctx: ctx,
		wg:  wg,
		api: api,

		msgCh:          make(chan *tgbotapi.MessageConfig, 200),
		userCollection: make(map[int64]*entities.User),

		commandManager: commandManager,
		messageManager: messageManager,

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
			if msg := b.proceedUpdate(update); msg != nil {
				b.msgCh <- msg
			}
		case <-b.ctx.Done():
			b.wg.Done()
			zap.S().Info("Telegram bot API stopped.")
			return
		}
	}
}

func (b *Bot) send() {
	for msg := range b.msgCh {
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

func (b *Bot) proceedUpdate(update tgbotapi.Update) *tgbotapi.MessageConfig {
	var msg *tgbotapi.MessageConfig
	var err error

	if update.Message != nil && update.Message.Chat != nil {
		user := b.prepareUser(update.Message)

		if !update.Message.IsCommand() {
			msg, err = b.messageManager.proceed(user, update)
		} else {
			msg, err = b.commandManager.proceed(user, update)
		}

		if err != nil {
			msg = b.proceedError(user, err)
		}
	}

	return msg
}

func (b *Bot) proceedError(user *entities.User, err error) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(
		user.ChatID,
		b.translator.Translate(err.Error(), "en", nil),
	)

	return &msg
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
