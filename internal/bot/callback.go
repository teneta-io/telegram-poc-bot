package bot

import (
	"encoding/json"
)

type CallbackData struct {
	MessageType string           `json:"t"`
	Data        *json.RawMessage `json:"d"`
}

type ProviderVCPUSelectorCallbackData struct {
	Count string `json:"c"`
}
