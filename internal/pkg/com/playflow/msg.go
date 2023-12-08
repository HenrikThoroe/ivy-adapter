package playflow

import (
	"encoding/json"
	"errors"
)

type Flow struct {
}

// NewFlow creates a new Flow.
func NewFlow() Flow {
	return Flow{}
}

// Parse parses a message from the server.
// The key is used to determine the type of the message.
func (f Flow) Parse(key string, data []byte) (any, error) {
	switch key {
	case "move-req-msg":
		return parse[MoveRequestMsg](data)
	case "update-msg":
		return parse[UpdateMsg](data)
	default:
		return nil, errors.New("invalid key")
	}
}

func parse[T any](data []byte) (T, error) {
	var m T
	err := json.Unmarshal(data, &m)
	return m, err
}

type MoveRequestMsg struct {
	Key     string   `json:"key"`
	History []string `json:"history"`
	Time    int      `json:"time"`
	Start   string   `json:"start"`
}

type UpdateMsg struct {
	Key     string   `json:"key"`
	History []string `json:"history"`
	Start   string   `json:"start"`
}
