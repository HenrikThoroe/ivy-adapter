package uci

import (
	"os"
	"strings"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
)

// OptionType is an enum for the type of an option.
type OptionType string

const (
	Spin   OptionType = "spin"   // Integer
	Check  OptionType = "check"  // Boolean
	Combo  OptionType = "combo"  // String, predefined values
	Button OptionType = "button" // No value
	String OptionType = "string" // String
)

// UCI is a wrapper for the communication between the adapter and the engine.
// It provides a simple interface to send commands to the engine and receive its responses.
type UCI struct {
	engine  *mgmt.Connection
	header  string
	options []OptionConfig
}

// ScoreType is an enum for the type of a score.
type ScoreType string

const (
	CP   ScoreType = "cp"   // Centipawns
	Mate ScoreType = "mate" // Mate in x moves
)

// Score is a wrapper for the score returned by the engine.
type Score struct {
	Type       ScoreType `json:"type"`       // Type of the score
	Value      int       `json:"value"`      // Value of the score
	Lowerbound bool      `json:"lowerbound"` // Whether the score is a lowerbound
	Upperbound bool      `json:"upperbound"` // Whether the score is an upperbound
}

// MoveInfo is a wrapper for the information returned by the engine after a move.
type MoveInfo struct {
	Move              string   `json:"move,omitempty"`              // The move itself
	Depth             int      `json:"depth,omitempty"`             // The depth the engine searched to
	SelDepth          int      `json:"selDepth,omitempty"`          // The selective depth the engine searched to
	Time              int      `json:"time,omitempty"`              // The time the engine searched in ms
	Nodes             int      `json:"nodes,omitempty"`             // The number of nodes the engine searched
	Pv                []string `json:"pv,omitempty"`                // The principal variation
	MultiPv           int      `json:"multipv,omitempty"`           // The multipv number
	Score             Score    `json:"score,omitempty"`             // The score of the move
	CurrentMove       string   `json:"currentMove,omitempty"`       // The current move the engine is searching
	CurrentMoveNumber int      `json:"currentMoveNumber,omitempty"` // The current move number
	HashFull          int      `json:"hashFull,omitempty"`          // The number of hash entries the engine searched
	Nps               int      `json:"nps,omitempty"`               // The number of nodes per second the engine searched
	TbHits            int      `json:"tbhits,omitempty"`            // The number of tablebase hits
	Sbhits            int      `json:"sbhits,omitempty"`            // The number of syzygy tablebase hits
	CpuLoad           int      `json:"cpuload,omitempty"`           // The CPU load in percent
	String            string   `json:"string,omitempty"`            // Additional information
	Refutation        []string `json:"refutation,omitempty"`        // The refutation to the current move
	Currline          []string `json:"currline,omitempty"`          // The current line the engine is searching
}

// Option contains the name and value of an option, which can be set.
type Option struct {
	Name  string
	Value string
}

// OptionConfig is parsed from the engine's response to the "uci" command.
// It contains information about the engine's options that can be set.
type OptionConfig struct {
	Name string
	Type OptionType
	Min  int
	Max  int
	Def  string
	Var  []string
}

// New returns a new UCI struct.
// It launches the engine and returns a connection to it.
func New(e *mgmt.EngineInstance) (*UCI, error) {
	conn, err := e.Launch()

	if err != nil {
		return nil, err
	}

	return &UCI{
		engine:  conn,
		options: make([]OptionConfig, 0),
	}, nil
}

// Close kills the engine process.
func (uci *UCI) Close() {
	proc, err := os.FindProcess(uci.engine.Pid)

	if err != nil {
		return
	}

	proc.Kill()
}

// GetOptionConfig returns the option configuration for the given option name.
// If the option does not exist, nil is returned.
// The option name is case insensitive.
func (u UCI) GetOptionConfig(name string) *OptionConfig {
	for _, opt := range u.options {
		if strings.EqualFold(opt.Name, name) {
			return &opt
		}
	}

	return nil
}

// String returns the string representation of an option as the command to send to the engine.
func (o Option) String() string {
	if o.Value == "" {
		return "option name " + o.Name
	}

	return "option name " + o.Name + " value " + o.Value
}

// Response returns the response to an option configuration.
// The response is an Option struct containing the name of the OptionConfig
// and the given value.
// If the OptionConfig is of type Button, the value is ignored.
func (o OptionConfig) Response(value string) Option {
	if o.Type == Button {
		return Option{
			Name: o.Name,
		}
	}

	return Option{
		Name:  o.Name,
		Value: value,
	}
}
