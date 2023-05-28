package testflow

import (
	"encoding/json"
	"errors"
)

type version_t struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

type time_t struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type options_t struct {
	HashSize int `json:"hash"`
	Threads  int `json:"threads"`
}

type engine_t struct {
	Name    string    `json:"name"`
	Version version_t `json:"version"`
	Time    time_t    `json:"timeControl"`
	Options options_t `json:"options"`
}

type suite_t struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Iterations int        `json:"iterations"`
	Engines    []engine_t `json:"engines"`
}

type Flow struct {
}

// NewFlow returns a new Flow.
func NewFlow() Flow {
	return Flow{}
}

// Parse parses the data into a message.
// The key is used to determine the type of the message.
func (f Flow) Parse(key string, data []byte) (any, error) {
	switch key {
	case "registered":
		return parse[RegisteredMsg](data)
	case "start":
		return parse[StartMsg](data)
	default:
		return nil, errors.New("invalid key")
	}
}

func parse[T any](data []byte) (T, error) {
	var m T
	err := json.Unmarshal(data, &m)
	return m, err
}

type RegisteredMsg struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type StartMsg struct {
	Key                  string  `json:"key"`
	Session              string  `json:"session"`
	Suite                suite_t `json:"suite"`
	RecommendedBatchSize int     `json:"recommendedBatchSize"`
}
