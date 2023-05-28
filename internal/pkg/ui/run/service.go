package run

import (
	"errors"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com/playflow"
)

func (m *model) loop() {
	reg := playflow.BuildRegisterCmd(m.game, m.color)

	m.gameData.client.Commands <- reg

	resp := <-m.gameData.client.Messages

	if info, ok := resp.(playflow.PlayerInfoMsg); ok {
		m.gameData.player = info.Player
		m.update <- true
	} else {
		m.gameData.err = errors.New("did not receive player information")
		m.update <- true
		return
	}

	for {
		msg := <-m.gameData.client.Messages

		switch msg := msg.(type) {
		case playflow.MoveRequestMsg:
			m.handleMoveRequest(msg)
		case playflow.GameStateMsg:
			m.gameData.moves = msg.Moves
			m.gameData.wtime = msg.Time.White
			m.gameData.btime = msg.Time.Black

			if msg.State != "active" {
				m.gameData.winner = msg.Winner
				m.gameData.reason = msg.Reason
				m.update <- true
				return
			}
		case playflow.ErrorMsg:
			m.gameData.err = errors.New("(Server Error) " + msg.Message)
		default:
			return
		}

		m.update <- true
	}
}

func (m *model) handleMoveRequest(msg playflow.MoveRequestMsg) {
	if msg.PlayerColor != m.color {
		m.gameData.err = errors.New("received move request for wrong color")
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
	m.gameData.client.Commands <- playflow.BuildMoveCmd(m.game, m.gameData.player, info.Move)
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
