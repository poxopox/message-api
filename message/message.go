package message

import (
	"encoding/json"
	"time"
)

type Message struct {
	User string `json:"user"`
	Message string `json:"message"`
	TimeStamp time.Time `json:"ts"`
}

func (msg *Message) ToJson () (string, error)  {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}