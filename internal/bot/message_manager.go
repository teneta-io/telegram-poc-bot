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
	tmp := strings.ReplaceAll(message, " ", "")
	ports := strings.Split(tmp, ",")

	errs := user.SetPorts(ports)
	user.State = nextState

	var messages = []string{}

	if len(errs) != len(ports) {
		messages = append(messages, m.translator.Translate("ports_limit_added", "en",
			map[string]interface{}{"ports": user.ProviderConfig.Ports.String()}))
	} else {
		if err := m.userService.Save(user); err != nil {
			return nil, err
		}
	}

	for port, err := range errs {
		args := map[string]interface{}{"availableProtocols": entities.AvailableProtocols, "port": port}
		msg, ok := m.translator.TryTranslate(err.Error(), "en", args)

		if !ok {
			msg = m.translator.Translate("unknown_ports_error", "en", args)
		}

		messages = append(messages, msg)
	}

	msg := tgbotapi.NewMessage(user.ChatID, strings.Join(messages, "\n"))

	return &msg, nil
}
