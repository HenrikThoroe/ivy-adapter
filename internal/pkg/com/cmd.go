package com

// Command is an interface for all commands that can be sent to the backend.
type Command interface {
	// Encode encodes the command into a JSON string.
	Encode() string
}
