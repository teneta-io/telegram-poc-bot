package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) actAsCustomerInlineKeyboard(language string) *tgbotapi.InlineKeyboardMarkup {
	keyboard := &tgbotapi.InlineKeyboardMarkup{}

	row := []tgbotapi.InlineKeyboardButton{
		{
			Text:         b.translator.Translate("create_task", language, nil),
			CallbackData: (&CallbackData{MessageType: CustomerCreateTaskMessageType}).ToString(),
		},
		{
			Text:         b.translator.Translate("task_list", language, nil),
			CallbackData: (&CallbackData{MessageType: CustomerTaskListMessageType}).ToString(),
		},
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	return keyboard
}
