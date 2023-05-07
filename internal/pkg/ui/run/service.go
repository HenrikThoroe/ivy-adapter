package run

import (
	"errors"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com"
)

func (m *model) loop() {
	reg := com.BuildRegisterCmd(m.game, m.color)

	m.gameData.client.Commands <- reg

	resp := <-m.gameData.client.Messages

	if info, ok := resp.(com.PlayerInfoMsg); ok {
		m.gameData.player = info.Player
		m.update <- true
	} else {
		m.gameData.err = errors.New("Did not receive player information")
		m.update <- true
		return
	}

	for {
		msg := <-m.gameData.client.Messages

		switch msg.(type) {
		case com.MoveRequestMsg:
			m.handleMoveRequest(msg.(com.MoveRequestMsg))
		case com.GameStateMsg:
			m.gameData.moves = msg.(com.GameStateMsg).Moves
			m.gameData.wtime = msg.(com.GameStateMsg).Time.White
			m.gameData.btime = msg.(com.GameStateMsg).Time.Black

			if msg.(com.GameStateMsg).State != "active" {
				m.gameData.winner = msg.(com.GameStateMsg).Winner
				m.gameData.reason = msg.(com.GameStateMsg).Reason
				m.update <- true
				return
			}
		case com.ErrorMsg:
			m.gameData.err = errors.New("(Server Error) " + msg.(com.ErrorMsg).Message)
		default:
			return
		}

		m.update <- true
	}
}

func (m *model) handleMoveRequest(msg com.MoveRequestMsg) {
	if msg.PlayerColor != m.color {
		m.gameData.err = errors.New("Received move request for wrong color")
		m.update <- true
		return
	}

	remaining := msg.Time.White

	if m.color == "black" {
		remaining = msg.Time.Black
	}

	m.gameData.ttm = m.getMoveTime(remaining, len(msg.Moves))
	m.gameData.engine.SetMoves(msg.Moves...)
	info := m.gameData.engine.GetMove(m.gameData.ttm)
	m.gameData.client.Commands <- com.BuildMoveCmd(m.game, m.gameData.player, info.Move)
}

func (m *model) getMoveTime(remaining int, moveCount int) int {
	blitz := 100
	quick := 300
	expectedMoves := 40

	if remaining < 3000 {
		return blitz
	}

	if remaining < 10000 {
		return quick
	}

	if moveCount/2 < expectedMoves {
		r := expectedMoves - moveCount/2
		return remaining / int(r)
	}

	return int(float64(remaining) * float64(0.02))
}
