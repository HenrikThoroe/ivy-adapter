package playflow

import (
	"encoding/json"
	"strings"
)

type UpdateReqCmd struct {
	Key string `json:"command"`
}

type CheckInCmd struct {
	Player string `json:"player"`
	Key    string `json:"key"`
}

type MoveCmd struct {
	Move string `json:"move"`
	Key  string `json:"key"`
}

func (c UpdateReqCmd) Encode() string {
	return encode(c)
}

func (c CheckInCmd) Encode() string {
	return encode(c)
}

func (c MoveCmd) Encode() string {
	return encode(c)
}

func BuildUpdateReqCmd(game string, player string) UpdateReqCmd {
	return UpdateReqCmd{
		Key: "update-req-msg",
	}
}

func BuildCheckInCmd(player string) CheckInCmd {
	return CheckInCmd{
		Player: player,
		Key:    "check-in-msg",
	}
}

func BuildMoveCmd(move string) MoveCmd {
	return MoveCmd{
		Move: move,
		Key:  "move-msg",
	}
}

func encode(data interface{}) string {
	buf := new(strings.Builder)
	enc := json.NewEncoder(buf)
	enc.Encode(data)

	return buf.String()
}
