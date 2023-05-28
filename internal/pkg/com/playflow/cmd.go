package playflow

import (
	"encoding/json"
	"strings"
)

// RegisterCmd is a command that registers a player for a game.
type RegisterCmd struct {
	Game  string `json:"game"`
	Color string `json:"color"`
	Key   string `json:"command"`
}

// GameStateCmd is a command that requests the current game state.
type GameStateCmd struct {
	Game   string `json:"game"`
	Player string `json:"player"`
	Key    string `json:"command"`
}

// CheckInCmd is a command that checks in a player for a game.
// Used when the game is already started and the player slot has been reserved.
// The slot must be currently empty. Can be used to reconnect when the
// initial connection was lost.
type CheckInCmd struct {
	Game   string `json:"game"`
	Player string `json:"player"`
	Key    string `json:"command"`
}

// MoveCmd is a command that sends a move to the backend.
type MoveCmd struct {
	Game   string `json:"game"`
	Player string `json:"player"`
	Move   string `json:"move"`
	Key    string `json:"command"`
}

// ResignCmd is a command that resigns the game.
type ResignCmd struct {
	Game   string `json:"game"`
	Player string `json:"player"`
	Key    string `json:"command"`
}

func (c RegisterCmd) Encode() string {
	return encode(c)
}

func (c GameStateCmd) Encode() string {
	return encode(c)
}

func (c CheckInCmd) Encode() string {
	return encode(c)
}

func (c MoveCmd) Encode() string {
	return encode(c)
}

func (c ResignCmd) Encode() string {
	return encode(c)
}

// BuildRegisterCmd builds a RegisterCmd with the given parameters.
func BuildRegisterCmd(game string, color string) RegisterCmd {
	return RegisterCmd{
		Game:  game,
		Color: color,
		Key:   "register",
	}
}

// BuildGameStateCmd builds a GameStateCmd with the given parameters.
func BuildGameStateCmd(game string, player string) GameStateCmd {
	return GameStateCmd{
		Game:   game,
		Player: player,
		Key:    "game-state",
	}
}

// BuildCheckInCmd builds a CheckInCmd with the given parameters.
func BuildCheckInCmd(game string, player string) CheckInCmd {
	return CheckInCmd{
		Game:   game,
		Player: player,
		Key:    "check-in",
	}
}

// BuildMoveCmd builds a MoveCmd with the given parameters.
func BuildMoveCmd(game string, player string, move string) MoveCmd {
	return MoveCmd{
		Game:   game,
		Player: player,
		Move:   move,
		Key:    "move",
	}

}

// BuildResignCmd builds a ResignCmd with the given parameters.
func BuildResignCmd(game string, player string) ResignCmd {
	return ResignCmd{
		Game:   game,
		Player: player,
		Key:    "resign",
	}
}

func encode(data interface{}) string {
	buf := new(strings.Builder)
	enc := json.NewEncoder(buf)
	enc.Encode(data)

	return buf.String()
}
