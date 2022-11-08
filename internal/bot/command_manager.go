package bot

import (
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"teneta-tg/internal/entities"
	"teneta-tg/internal/translator"
)

const (
	StartCommand = "start"

	ActAsProviderCommand = "act_as_provider"
	ActAsCustomerCommand = "act_as_customer"

	AddVCPULimitCommand    = "vcpu"
	AddRAMLimitCommand     = "ram"
	AddStorageLimitCommand = "storage"
	AddNetworkLimitCommand = "network"

	OpenPortCommand = "ports"

	AboutCommand = "about"
)

type CommandManager struct {
	keyboardManager *KeyboardManager
	translator      *translator.Translator
}

func NewCommandManager(keyboardManager *KeyboardManager, translator *translator.Translator) *CommandManager {
	return &CommandManager{
		keyboardManager: keyboardManager,
		translator:      translator,
	}
}

func (m *CommandManager) proceed(user *entities.User, update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	switch update.Message.Command() {
	case StartCommand:
		return m.execStartCommand(user)
	case ActAsProviderCommand:
		return m.execActAsProviderCommand(user)
	case ActAsCustomerCommand:

	case AddVCPULimitCommand:
		return m.execAddResourceCommand(user, "vcpu", AddVCPULimitState)
	case AddRAMLimitCommand:
		return m.execAddResourceCommand(user, "ram", AddRamLimitState)
	case AddStorageLimitCommand:
		return m.execAddResourceCommand(user, "storage", AddStorageLimitState)
	case AddNetworkLimitCommand:
		return m.execAddResourceCommand(user, "network", AddNetworkLimitState)
	case OpenPortCommand:
		return m.execOpenPortCommand(user, AddPortsState)
	case AboutCommand:
	}

	return nil, errors.New("undefined_command")
}

func (m *CommandManager) execStartCommand(user *entities.User) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(
		user.ChatID,
		m.translator.Translate("start_command_response", "en", nil),
	)

	return &msg, nil
}

func (m *CommandManager) execActAsProviderCommand(user *entities.User) (*tgbotapi.MessageConfig, error) {
	if user.ProviderConfig == nil {
		user.ProviderConfig = &entities.Provider{}
	}

	msg := tgbotapi.NewMessage(
		user.ChatID,
		m.translator.Translate("act_as_provider_response", "en", nil),
	)

	return &msg, nil
}

func (m *CommandManager) execAddResourceCommand(user *entities.User, t string, state int) (*tgbotapi.MessageConfig, error) {
	if !user.IsProvider() {
		return nil, errors.New("user_not_registered_as_provider")
	}

	user.State = state
	msg := tgbotapi.NewMessage(
		user.ChatID,
		m.translator.Translate(fmt.Sprintf("%s_limit_add_start", t), "en", nil),
	)

	return &msg, nil
}

func (m *CommandManager) execOpenPortCommand(user *entities.User, state int) (*tgbotapi.MessageConfig, error) {
	if !user.IsProvider() {
		return nil, errors.New("user_not_registered_as_provider")
	}

	user.State = state
	msg := tgbotapi.NewMessage(
		user.ChatID,
		m.translator.Translate("open_ports_start", "en", nil),
	)

	return &msg, nil
}
