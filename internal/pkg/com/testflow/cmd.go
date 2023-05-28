package testflow

import (
	"encoding/json"
	"strings"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/sys"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/uci"
)

type GameMoveHistory [][]uci.MoveInfo

type ReportCmd struct {
	Key     string            `json:"command"`
	Session string            `json:"session"`
	Moves   []GameMoveHistory `json:"moves"`
}

type RegisterCmd struct {
	Key      string     `json:"command"`
	Name     string     `json:"name"`
	DeviceId string     `json:"deviceId"`
	Hardware sys.Device `json:"hardware"`
}

func (c RegisterCmd) Encode() string {
	return encode(c)
}

func (c ReportCmd) Encode() string {
	return encode(c)
}

func BuildRegisterCmd() RegisterCmd {
	hardware, identifier := sys.DeviceInfo()

	return RegisterCmd{
		Key:      "register",
		Name:     identifier.Name,
		DeviceId: identifier.ID,
		Hardware: hardware,
	}
}

func BuildReportCmd(session string, moves []GameMoveHistory) ReportCmd {
	return ReportCmd{
		Key:     "report",
		Session: session,
		Moves:   moves,
	}
}

func encode(data interface{}) string {
	buf := new(strings.Builder)
	enc := json.NewEncoder(buf)
	enc.Encode(data)

	return buf.String()
}
