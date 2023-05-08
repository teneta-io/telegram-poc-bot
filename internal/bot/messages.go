package bot

import (
	"fmt"
	"strconv"
	"strings"
	"teneta-tg/internal/entities"
)

const (
	resourceTypeVCPU    = "vcpu"
	resourceTypeRam     = "ram"
	resourceTypeStorage = "storage"
	resourceTypeNetwork = "network"
	resourceTypePorts   = "ports"
)

var (
	resourceStateType = map[int]string{
		addVCPULimitState:    resourceTypeVCPU,
		addRamLimitState:     resourceTypeRam,
		addStorageLimitState: resourceTypeStorage,
		addNetworkLimitState: resourceTypeNetwork,
		addPortsState:        resourceTypePorts,
	}
)

func (b *Bot) proceedMessage(user *entities.User, text string) {
	switch user.State {
	case addVCPULimitState, addRamLimitState, addStorageLimitState, addNetworkLimitState:
		b.proceedAddResourceLimit(user, resourceStateType[user.State], text, actAsProviderState)
		b.userService.Save(user)

		return
	case addPortsState:
		b.proceedAddPort(user, resourceTypePorts, text, actAsProviderState)
		b.userService.Save(user)

		return
	}

	b.messageCh <- &MessageResponse{
		ChatId: user.ChatID,
		Text:   b.translator.Translate("something_wrong", "en", nil),
	}
}

func (b *Bot) proceedAddResourceLimit(user *entities.User, t, message string, nextState int) {
	v, err := strconv.ParseInt(message, 10, 64)

	if err != nil {
		b.messageCh <- &MessageResponse{
			ChatId: user.ChatID,
			Text:   b.translator.Translate(fmt.Sprintf("invalid_%s_value_error", t), "en", nil),
		}

		return
	}

	switch t {
	case resourceTypeVCPU:
		user.ProviderConfig.VCPU = v
	case resourceTypeRam:
		user.ProviderConfig.Ram = v
	case resourceTypeStorage:
		user.ProviderConfig.Storage = v
	case resourceTypeNetwork:
		user.ProviderConfig.Network = v
	}

	user.State = nextState

	b.messageCh <- &MessageResponse{
		ChatId: user.ChatID,
		Text:   b.translator.Translate(fmt.Sprintf("%s_added", t), "en", map[string]interface{}{"count": v}),
	}

	return
}

func (b *Bot) proceedAddPort(user *entities.User, t, message string, nextState int) {
	tmp := strings.ReplaceAll(message, " ", "")
	ports := strings.Split(tmp, ",")

	errs := user.SetPorts(ports)
	user.State = nextState

	var messages = []string{}

	if len(errs) != len(ports) {
		messages = append(messages, b.translator.Translate("ports_added", "en",
			map[string]interface{}{"ports": user.ProviderConfig.Ports.String()}))
	}

	for port, err := range errs {
		args := map[string]interface{}{"availableProtocols": entities.AvailableProtocols, "port": port}
		msg, ok := b.translator.TryTranslate(err.Error(), "en", args)

		if !ok {
			msg = b.translator.Translate("unknown_ports_error", "en", args)
		}

		messages = append(messages, msg)
	}

	b.messageCh <- &MessageResponse{
		ChatId: user.ChatID,
		Text:   strings.Join(messages, "\n"),
	}

	return
}
