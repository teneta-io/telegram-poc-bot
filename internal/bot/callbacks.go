package bot

import (
	"encoding/json"
	"go.uber.org/zap"
	"teneta-tg/internal/entities"
)

const (
	CustomerCreateTaskMessageType = "cus_ct"
	CustomerTaskListMessageType   = "cus_tl"
)

type CallbackData struct {
	MessageType string           `json:"t"`
	Data        *json.RawMessage `json:"d"`
}

func (d *CallbackData) ToString() *string {
	callbackDataJson, err := json.Marshal(d)
	if err != nil {
		zap.S().Error(err)

		return nil
	}

	str := string(callbackDataJson)

	return &str
}

func (b *Bot) proceedCallback(user *entities.User, callbackData string) {
	var callback *CallbackData

	if err := json.Unmarshal([]byte(callbackData), &callback); err != nil {
		zap.S().Error(err)
		b.response(user, "something_wrong", nil, nil, nil)

		return
	}

	switch callback.MessageType {
	case CustomerCreateTaskMessageType:
		zap.S().Info("create task")
	case CustomerTaskListMessageType:
		zap.S().Info("task list")
	}
}
