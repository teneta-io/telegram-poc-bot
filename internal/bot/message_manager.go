package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"teneta-tg/internal/entities"
	"teneta-tg/internal/translator"
)

type MessageManager struct {
	translator *translator.Translator
}

func NewMessageManager(translator *translator.Translator) *MessageManager {
	return &MessageManager{
		translator: translator,
	}
}

func (m *MessageManager) proceed(user *entities.User, update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	switch user.State {
	case AddVCPULimitState:
		return m.proceedAddResourceLimit(user, "vcpu", update.Message.Text, DefaultState)
	case AddRamLimitState:
		return m.proceedAddResourceLimit(user, "ram", update.Message.Text, DefaultState)
	case AddStorageLimitState:
		return m.proceedAddResourceLimit(user, "storage", update.Message.Text, DefaultState)
	case AddNetworkLimitState:
		return m.proceedAddResourceLimit(user, "network", update.Message.Text, DefaultState)
	}

	return nil, nil
}

func (m *MessageManager) proceedAddResourceLimit(user *entities.User, t, message string, nextState int) (*tgbotapi.MessageConfig, error) {
	v, err := strconv.ParseInt(message, 10, 64)

	if err != nil {
		return nil, fmt.Errorf("invalid_%s_value_error", t)
	}

	user.ProviderConfig.VCPU = v
	user.State = nextState
	msg := tgbotapi.NewMessage(
		user.ChatID,
		m.translator.Translate(fmt.Sprintf("%s_limit_added", t), "en", map[string]interface{}{"count": v}),
	)

	return &msg, nil
}
