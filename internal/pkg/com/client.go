package com

import (
	"encoding/json"
	"errors"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
	"github.com/gorilla/websocket"
)

// Client represents a connection to the game management server.
type Client struct {
	// Messages is a channel that receives messages from the server and forwards the parsed message.
	Messages chan any

	// Commands is a channel that receives commands from the client and sends the encoded command over the ws connection.
	Commands chan Command

	// Errors is a channel that distributes any errors that occur during the connection.
	Errors chan error
}

// Connect establishes a connection to the game management server based on the configuration in the environment.
func Connect() (*Client, error) {
	msgChan := make(chan any)
	cmdChan := make(chan Command)
	errChan := make(chan error)
	conn, _, err := websocket.DefaultDialer.Dial(conf.GetGameServerConfig().GetURL(), nil)

	if err != nil {
		return nil, err
	}

	go handleReceive(msgChan, errChan, conn)
	go handleSend(cmdChan, errChan, conn)

	return &Client{
		Messages: msgChan,
		Commands: cmdChan,
		Errors:   errChan,
	}, nil
}

func handleSend(cmdChan chan Command, errChan chan error, conn *websocket.Conn) {
	for {
		cmd := <-cmdChan
		message := cmd.Encode()

		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			errChan <- err
			continue
		}
	}
}

func handleReceive(msgChan chan any, errChan chan error, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			errChan <- err
			continue
		}

		data := map[string]any{}

		if err := json.Unmarshal(message, &data); err != nil {
			errChan <- err
			continue
		}

		key := data["key"]
		var e error
		var d any

		switch key {
		case "move-request":
			d, e = parseMoveRequestMsg(string(message))
		case "player-info":
			d, e = parsePlayerInfoMsg(string(message))
		case "game-state":
			d, e = parseGameStateMsg(string(message))
		case "error":
			d, e = parseErrorMsg(string(message))
		case "pong":
			d, e = parsePongMsg(string(message))
		default:
			e = errors.New("Received message with unknown key attribute")
		}

		if e != nil {
			errChan <- e
			continue
		}

		msgChan <- d
	}
}
