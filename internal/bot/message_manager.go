package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"teneta-tg/internal/entities"
	"teneta-tg/internal/services"
	"teneta-tg/internal/translator"
)

type MessageManager struct {
	translator  *translator.Translator
	userService services.UserService
}

func NewMessageManager(translator *translator.Translator, userService services.UserService) *MessageManager {
	return &MessageManager{
		translator:  translator,
		userService: userService,
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
	case AddPortsState:
		return m.proceedAddPort(user, "ports", update.Message.Text, DefaultState)
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

func (m *MessageManager) proceedAddPort(user *entities.User, t, message string, nextState int) (*tgbotapi.MessageConfig, error) {
	tmp := strings.ReplaceAll(message, "[", "")
	tmp = strings.ReplaceAll(tmp, "]", "")
	tmp = strings.ReplaceAll(tmp, " ", "")

	err := user.SetPorts(strings.Split(tmp, ","))
	user.State = nextState

	var msg tgbotapi.MessageConfig

	switch err {
	case entities.ErrWrongPortFormat:
		msg = tgbotapi.NewMessage(
			user.ChatID,
			m.translator.Translate("invalid_ports_format", "en", nil),
		)
	case entities.ErrNumberCanNotBePort:
		msg = tgbotapi.NewMessage(
			user.ChatID,
			m.translator.Translate("invalid_ports_number_error", "en", nil),
		)
	case entities.ErrProtocolIsNotSupported:
		msg = tgbotapi.NewMessage(
			user.ChatID,
			m.translator.Translate("invalid_ports_protocol_error", "en",
				map[string]interface{}{"protocols": entities.AvailableProtocols}),
		)
	case nil:
		msg = tgbotapi.NewMessage(
			user.ChatID,
			m.translator.Translate("ports_limit_added", "en",
				map[string]interface{}{"ports": user.ProviderConfig.Ports.String()}),
		)

		if err := m.userService.Save(user); err != nil {
			return nil, err
		}

	}

	return &msg, nil
}
