package uci

import (
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
)

// UCI is a wrapper for the communication between the adapter and the engine.
// It provides a simple interface to send commands to the engine and receive its responses.
type UCI struct {
	engine *mgmt.Connection
	header string
}

// ScoreType is an enum for the type of a score.
type ScoreType int

const (
	CP   ScoreType = iota // Centipawns
	Mate                  // Mate in x moves
)

// Score is a wrapper for the score returned by the engine.
type Score struct {
	Type       ScoreType // Type of the score
	Value      int       // Value of the score
	Lowerbound bool      // Whether the score is a lowerbound
	Upperbound bool      // Whether the score is an upperbound
}

// MoveInfo is a wrapper for the information returned by the engine after a move.
type MoveInfo struct {
	Move              string   // The move itself
	Depth             int      // The depth the engine searched to
	SelDepth          int      // The selective depth the engine searched to
	Time              int      // The time the engine searched in ms
	Nodes             int      // The number of nodes the engine searched
	Pv                []string // The principal variation
	MultiPv           int      // The multipv number
	Score             Score    // The score of the move
	CurrentMove       string   // The current move the engine is searching
	CurrentMoveNumber int      // The current move number
	HashFull          int      // The number of hash entries the engine searched
	Nps               int      // The number of nodes per second the engine searched
	TbHits            int      // The number of tablebase hits
	Sbhits            int      // The number of syzygy tablebase hits
	CpuLoad           int      // The CPU load in percent
	String            string   // Additional information
	Refutation        []string // The refutation to the current move
	Currline          []string // The current line the engine is searching
}

// New returns a new UCI struct.
// It launches the engine and returns a connection to it.
func New(e *mgmt.EngineInstance) (*UCI, error) {
	conn, err := e.Launch()

	if err != nil {
		return nil, err
	}

	return &UCI{
		engine: conn,
	}, nil
}
