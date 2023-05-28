package playflow

import (
	"encoding/json"
	"errors"
)

type time_t struct {
	White int `json:"white"`
	Black int `json:"black"`
}

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
	case "move-request":
		return parse[MoveRequestMsg](data)
	case "player-info":
		return parse[PlayerInfoMsg](data)
	case "game-state":
		return parse[GameStateMsg](data)
	case "error":
		return parse[ErrorMsg](data)
	case "pong":
		return parse[PongMsg](data)
	default:
		return nil, errors.New("invalid key")
	}
}

func parse[T any](data []byte) (T, error) {
	var m T
	err := json.Unmarshal(data, &m)
	return m, err
}

// MoveRequestMsg is a message sent by the server to request a move from the
// client.
type MoveRequestMsg struct {
	Key         string   `json:"key"`
	PlayerColor string   `json:"playerColor"`
	Moves       []string `json:"moves"`
	Time        time_t   `json:"time"`
}

// PlayerInfoMsg is a message sent by the server to inform the client about
// the player's color and id.
type PlayerInfoMsg struct {
	Key    string `json:"key"`
	Player string `json:"player"`
	Color  string `json:"color"`
	Game   string `json:"game"`
}

// GameStateMsg is a message sent by the server to inform the client about
// the current state of the game.
type GameStateMsg struct {
	Key         string   `json:"key"`
	PlayerColor string   `json:"playerColor"`
	Moves       []string `json:"moves"`
	Time        time_t   `json:"time"`
	State       string   `json:"state"`
	Winner      string   `json:"winner"`
	Reason      string   `json:"reason"`
}

// ErrorMsg is a message sent by the server to inform the client about an
// error that occurred.
type ErrorMsg struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

// PongMsg is a message sent by the client to the server to respond to a ping.
type PongMsg struct {
	Key string `json:"key"`
}
