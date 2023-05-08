package bot

import (
	"fmt"
	"teneta-tg/internal/entities"
)

const (
	startCommand = "start"

	actAsProviderCommand = "act_as_provider"
	actAsCustomerCommand = "act_as_customer"

	addVCPULimitCommand    = "vcpu"
	addRAMLimitCommand     = "ram"
	addStorageLimitCommand = "storage"
	addNetworkLimitCommand = "network"
	addPortCommand         = "ports"

	aboutCommand = "about"
)

var (
	resourceCommandState = map[string]int{
		addVCPULimitCommand:    addVCPULimitState,
		addRAMLimitCommand:     addRamLimitState,
		addStorageLimitCommand: addStorageLimitState,
		addNetworkLimitCommand: addNetworkLimitState,
		addPortCommand:         addPortsState,
	}
)

func (b *Bot) proceedCommand(user *entities.User, command string) {
	switch command {

	// SYSTEM
	case startCommand:
		b.proceedStartCommand(user)
		return
	case actAsProviderCommand:
		b.execActAsProviderCommand(user)
		b.userService.Save(user)
		return
	case actAsCustomerCommand:

	// PROVIDER
	case addVCPULimitCommand, addRAMLimitCommand, addStorageLimitCommand, addNetworkLimitCommand, addPortCommand:
		b.proceedAddResourceCommand(user, command, resourceCommandState[command])
		return

		// CUSTOMER
	}

	b.messageCh <- &MessageResponse{
		ChatId: user.ChatID,
		Text:   b.translator.Translate("undefined_command", "en", nil),
	}
}

func (b *Bot) proceedStartCommand(user *entities.User) {
	b.response(user, "start_command_response", nil)
}

func (b *Bot) execActAsProviderCommand(user *entities.User) {
	if user.ProviderConfig == nil {
		user.ProviderConfig = &entities.Provider{}
		user.ProviderConfig.ChatID = user.ChatID
	}

	user.State = actAsProviderState

	b.response(user, "act_as_provider_response", nil)
}

func (b *Bot) proceedActAsCustomerCommand(user *entities.User) {
	user.State = actAsCustomerState

	b.response(user, "act_as_customer_response", nil)
}

func (b *Bot) proceedAddResourceCommand(user *entities.User, t string, state int) {
	if !user.IsProvider() {
		b.response(user, "user_not_registered_as_provider", nil)

		return
	}

	if user.State != actAsProviderState {
		b.response(user, "user_current_context_is_not_provider", nil)

		return
	}

	user.State = state

	b.response(user, fmt.Sprintf("%s_add_start", t), nil)

	return
}
