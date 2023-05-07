package com

import (
	"encoding/json"
	"errors"
)

type time_t struct {
	White int `json:"white"`
	Black int `json:"black"`
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

func parsePongMsg(msg string) (PongMsg, error) {
	var m PongMsg
	err := json.Unmarshal([]byte(msg), &m)

	if m.Key != "pong" {
		return m, errors.New("Message has invalid key attribute")
	}

	return m, err
}

func parseMoveRequestMsg(msg string) (MoveRequestMsg, error) {
	var m MoveRequestMsg
	err := json.Unmarshal([]byte(msg), &m)

	if m.Key != "move-request" {
		return m, errors.New("Message has invalid key attribute")
	}

	return m, err
}

func parsePlayerInfoMsg(msg string) (PlayerInfoMsg, error) {
	var m PlayerInfoMsg
	err := json.Unmarshal([]byte(msg), &m)

	if m.Key != "player-info" {
		return m, errors.New("Message has invalid key attribute")
	}

	return m, err
}

func parseGameStateMsg(msg string) (GameStateMsg, error) {
	var m GameStateMsg
	err := json.Unmarshal([]byte(msg), &m)

	if m.Key != "game-state" {
		return m, errors.New("Message has invalid key attribute")
	}

	return m, err
}

func parseErrorMsg(msg string) (ErrorMsg, error) {
	var m ErrorMsg
	err := json.Unmarshal([]byte(msg), &m)

	if m.Key != "error" {
		return m, errors.New("Message has invalid key attribute")
	}

	return m, err
}
