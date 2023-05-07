package com

import (
	"encoding/json"
	"errors"
	"time"

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

	// ping is the last measured ping to the server.
	ping int64
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

	client := Client{
		Messages: msgChan,
		Commands: cmdChan,
		Errors:   errChan,
		ping:     -1,
	}

	go handleReceive(msgChan, errChan, conn)
	go handleSend(cmdChan, errChan, conn, &client.ping)

	return &client, nil
}

// Ping returns the last measured ping to the server.
// If the ping is not yet measured, it will be measured and returned.
func (c Client) Ping() int64 {
	if c.ping < 0 {
		start := time.Now().UnixMilli()
		c.Commands <- BuildPingCmd()
		resp := <-c.Messages

		if _, ok := resp.(PongMsg); !ok {
			return -1
		}

		end := time.Now().UnixMilli()

		return end - start
	}

	return c.ping
}

func handleSend(cmdChan chan Command, errChan chan error, conn *websocket.Conn, ping *int64) {
	for {
		cmd := <-cmdChan
		message := cmd.Encode()
		start := time.Now().UnixMilli()

		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			errChan <- err
			continue
		}

		end := time.Now().UnixMilli()
		*ping = end - start
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
