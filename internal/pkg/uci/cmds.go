package uci

import (
	"fmt"
	"strings"
)

// Setup sends the uci command to the engine and waits for the uciok response.
// It also parses the response for option configurations and saves them in
// the UCI struct.
func (u *UCI) Setup() {
	u.header = u.engine.Line()
	u.engine.Send("uci")

	u.engine.Read(func(line string) bool {
		if strings.HasPrefix(line, "option") {
			opt, e := parseOptionConfigStr(line)

			if e == nil {
				u.options = append(u.options, opt)
			}
		}

		return line == "uciok"
	})
}

// Start sends the ucinewgame command to the engine.
func (u *UCI) Start() {
	u.engine.Send("ucinewgame")
}

// StopSearch sends the stop command to the engine.
func (u *UCI) StopSearch() {
	u.engine.Send("stop")
}

// Quit sends the quit command to the engine.
func (u *UCI) Quit() {
	u.engine.Send("quit")
}

// SetPosition sends the position command to the engine with fen being
// the starting position and moves being the moves to be played.
func (u *UCI) SetPosition(fen string, moves ...string) {
	if len(moves) == 0 {
		u.engine.Send("position fen " + fen)
	} else {
		u.engine.Send("position fen " + fen + " moves " + strings.Join(moves, " "))
	}
}

// SetMoves is equal to SetPosition with fen being the standard starting position.
func (u *UCI) SetMoves(moves ...string) {
	if len(moves) == 0 {
		u.engine.Send("position startpos")
	} else {
		u.engine.Send("position startpos moves " + strings.Join(moves, " "))
	}
}

// GetMove sends the go command to the engine with movetime being the time
// the engine has to think about the next move. The function returns a pointer
// to a MoveInfo struct containing the information about the move.
// The function blocks until the engine has found a move.
func (u *UCI) GetMove(ms int) *MoveInfo {
	cmd := "go movetime " + fmt.Sprintf("%d", ms)
	var info *MoveInfo
	var last string

	u.engine.Scan(cmd, func(line string) bool {
		if strings.HasPrefix(line, "bestmove") {
			move := strings.Split(line, " ")[1]
			info = parseInfoStr(last)
			info.Move = move
			return true
		}

		last = line
		return false
	})

	return info
}

// IsEngineReady sends the isready command to the engine and returns true if
// the engine is ready.
func (u *UCI) IsEngineReady() bool {
	u.engine.Send("isready")
	resp := u.engine.Line()

	return resp == "readyok"
}

// SetOption sends the setoption command to the engine with the given option.
func (u *UCI) SetOption(option Option) {
	u.engine.Send("setoption name " + option.Name + " value " + option.Value)
}
