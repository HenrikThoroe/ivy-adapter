package com

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

// Flow is an interface that is used to parse messages sent by the server.
// Differnt flows are used for different types of messages.
type Flow interface {
	// Parse parses a message from the server.
	Parse(key string, data []byte) (any, error)
}

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

	// flow is the flow that is used to parse messages.
	flow Flow
}

// Connect establishes a connection to the game management server based on the configuration in the environment.
func Connect(url string, flow Flow) (*Client, error) {
	msgChan := make(chan any, 2)
	cmdChan := make(chan Command)
	errChan := make(chan error)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}

	client := Client{
		Messages: msgChan,
		Commands: cmdChan,
		Errors:   errChan,
		ping:     -1,
		flow:     flow,
	}

	go handleReceive(msgChan, errChan, conn, flow)
	go handleSend(cmdChan, errChan, conn, &client.ping)

	return &client, nil
}

// Ping returns the last measured ping to the server.
// If the ping is not yet measured, 0 is returned.
func (c Client) Ping() int64 {
	if c.ping < 0 {
		return 0
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

func handleReceive(msgChan chan any, errChan chan error, conn *websocket.Conn, flow Flow) {
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
		d, e := flow.Parse(key.(string), message)

		if e != nil {
			errChan <- e
			continue
		}

		msgChan <- d
	}
}
