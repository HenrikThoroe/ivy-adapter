package com

import (
	"reflect"
	"testing"
)

type mvreq_io struct {
	msg   string
	value MoveRequestMsg
}

var mvreq_io_table = []mvreq_io{
	{
		`{"key":"move-request","playerColor":"white","moves":["e2e4","e7e5"],"time":{"white":300,"black":300}}`,
		MoveRequestMsg{
			Key:         "move-request",
			PlayerColor: "white",
			Moves:       []string{"e2e4", "e7e5"},
			Time: time{
				White: 300,
				Black: 300,
			},
		},
	},
	{
		`{"key":"move-request","playerColor":"white","moves":[],"time":{"white":300,"black":300}}`,
		MoveRequestMsg{
			Key:         "move-request",
			PlayerColor: "white",
			Moves:       []string{},
			Time: time{
				White: 300,
				Black: 300,
			},
		},
	},
	{
		`{"key":"move-request","playerColor":"white","moves":[],"time":{"white":5400000,"black":5400000}}`,
		MoveRequestMsg{
			Key:         "move-request",
			PlayerColor: "white",
			Moves:       []string{},
			Time: time{
				White: 5400000,
				Black: 5400000,
			},
		},
	},
}

func TestParseMoveRequest(t *testing.T) {
	for _, io := range mvreq_io_table {
		msg, err := parseMoveRequestMsg(io.msg)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(msg, io.value) {
			t.Errorf("Expected %v, got %v", io.value, msg)
		}
	}
}
