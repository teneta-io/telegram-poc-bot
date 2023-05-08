package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageResponse struct {
	ChatId         int64
	Text           string
	InlineMarkup   *tgbotapi.InlineKeyboardMarkup
	ReplyMarkup    *tgbotapi.ReplyKeyboardMarkup
	ErrorMessage   bool
	IgnoreBusyLock bool
}
