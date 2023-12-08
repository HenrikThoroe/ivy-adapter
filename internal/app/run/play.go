package run

import (
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com/playflow"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/uci"
)

// Play starts a game against the game server.
// It will connect to the server and send a check-in message.
// After that it will wait for a move request and send the move to the server.
// This process will repeat until the server closes the connection.
func Play(ifc *uci.UCI, id string) error {
	flow := playflow.NewFlow()
	client, err := com.Connect(conf.GetGameServerConfig().GetURL(), flow)

	if err != nil {
		return err
	}

	checkIn := playflow.BuildCheckInCmd(id)
	closeChan := make(chan bool)

	client.Commands <- checkIn
	ifc.Setup()
	ifc.Start()

	go listenForErrors(client, closeChan)

	for {
		select {
		case <-closeChan:
			return nil
		case m := <-client.Messages:
			switch msg := m.(type) {
			case playflow.MoveRequestMsg:
				client.Commands <- playflow.BuildMoveCmd(fetchMove(ifc, msg.Time, msg.Start, msg.History))
			default:
				continue
			}
		}

	}
}

func listenForErrors(client *com.Client, signal chan bool) {
	for range client.Errors {
		client.Close()
		signal <- true
		break
	}
}

func fetchMove(ifc *uci.UCI, time int, start string, moves []string) string {
	ifc.SetPosition(start, moves...)
	move := ifc.GetMove(time)
	return move.Move
}
