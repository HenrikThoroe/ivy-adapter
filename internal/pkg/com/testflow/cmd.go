package testflow

import (
	"encoding/json"
	"strings"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/sys"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/uci"
)

// GameMoveHistory is a slice of moves for a single game.
type GameMoveHistory [][]uci.MoveInfo

// LogEntry is a struct that represents a single log entry.
type LogEntry struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Log is a slice of log entries.
type Log [][]LogEntry

// ReportCmd is a struct that represents a report command.
// This command is used to report the results of a batch of games.
type ReportCmd struct {
	Key     string            `json:"command"`
	Session string            `json:"session"`
	Moves   []GameMoveHistory `json:"moves"`
	Logs    []Log             `json:"logs"`
}

// RegisterCmd is a struct that represents a register command.
// This command is used to register a test driver.
type RegisterCmd struct {
	Key      string     `json:"command"`
	Name     string     `json:"name"`
	DeviceId string     `json:"deviceId"`
	Hardware sys.Device `json:"hardware"`
}

// Encode returns a string representation of the command.
func (c RegisterCmd) Encode() string {
	return encode(c)
}

// Encode returns a string representation of the command.
func (c ReportCmd) Encode() string {
	return encode(c)
}

// BuildRegisterCmd returns a RegisterCmd.
// It automatically fills in the name, device id and hardware fields.
func BuildRegisterCmd() RegisterCmd {
	hardware, identifier := sys.DeviceInfo()

	return RegisterCmd{
		Key:      "register",
		Name:     identifier.Name,
		DeviceId: identifier.ID,
		Hardware: hardware,
	}
}

// BuildReportCmd returns a ReportCmd with the given parameters.
func BuildReportCmd(session string, moves []GameMoveHistory, logs []Log) ReportCmd {
	return ReportCmd{
		Key:     "report",
		Session: session,
		Moves:   moves,
		Logs:    logs,
	}
}

func encode(data interface{}) string {
	buf := new(strings.Builder)
	enc := json.NewEncoder(buf)
	enc.Encode(data)

	return buf.String()
}
