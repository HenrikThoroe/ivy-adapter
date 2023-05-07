package uci

import (
	"reflect"
	"strings"
	"testing"
)

type score_io struct {
	in     string
	out    Score
	tokens int
}

type infostr_io struct {
	in  string
	out MoveInfo
}

var scores = []score_io{
	{"score cp 100", Score{Type: CP, Value: 100, Lowerbound: false, Upperbound: false}, 3},
	{"score cp 100 lowerbound", Score{Type: CP, Value: 100, Lowerbound: true, Upperbound: false}, 4},
	{"score cp 100 upperbound", Score{Type: CP, Value: 100, Lowerbound: false, Upperbound: true}, 4},
	{"score mate 100", Score{Type: Mate, Value: 100, Lowerbound: false, Upperbound: false}, 3},
	{"score mate 100 lowerbound", Score{Type: Mate, Value: 100, Lowerbound: true, Upperbound: false}, 4},
	{"score mate 100 upperbound", Score{Type: Mate, Value: 100, Lowerbound: false, Upperbound: true}, 4},
	{"score cp -100", Score{Type: CP, Value: -100, Lowerbound: false, Upperbound: false}, 3},
	{"score cp -100 lowerbound", Score{Type: CP, Value: -100, Lowerbound: true, Upperbound: false}, 4},
	{"score cp -100 upperbound", Score{Type: CP, Value: -100, Lowerbound: false, Upperbound: true}, 4},
	{"score mate -100", Score{Type: Mate, Value: -100, Lowerbound: false, Upperbound: false}, 3},
	{"score mate -100 lowerbound", Score{Type: Mate, Value: -100, Lowerbound: true, Upperbound: false}, 4},
	{"score mate -100 upperbound", Score{Type: Mate, Value: -100, Lowerbound: false, Upperbound: true}, 4},

	{"score cp 100 suffix", Score{Type: CP, Value: 100, Lowerbound: false, Upperbound: false}, 3},
	{"score cp 100 lowerbound suffix", Score{Type: CP, Value: 100, Lowerbound: true, Upperbound: false}, 4},
}

var infostrs = []infostr_io{
	{"info depth 1", MoveInfo{Depth: 1}},
	{"info depth 1 seldepth 2", MoveInfo{Depth: 1, SelDepth: 2}},
	{"info depth 1 seldepth 2 multipv 3", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 lowerbound", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100, Lowerbound: true}}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 upperbound", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100, Upperbound: true}}},
	{"info depth 1 seldepth 2 multipv 3 score mate 100", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: Mate, Value: 100}}},
	{"info depth 1 seldepth 2 multipv 3 score mate 100 lowerbound", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: Mate, Value: 100, Lowerbound: true}}},
	{"info depth 1 seldepth 2 multipv 3 score mate 100 upperbound", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: Mate, Value: 100, Upperbound: true}}},
	{"info depth 1 seldepth 2 multipv 3 score cp -100", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: -100}}},
	{"info depth 1 seldepth 2 multipv 3 score cp -100 lowerbound", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: -100, Lowerbound: true}}},
	{"info depth 1 seldepth 2 multipv 3 score cp -100 upperbound", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: -100, Upperbound: true}}},
	{"info depth 1 seldepth 2 multipv 3 score mate -100", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: Mate, Value: -100}}},
	{"info depth 1 seldepth 2 multipv 3 score mate -100 lowerbound", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: Mate, Value: -100, Lowerbound: true}}},
	{"info depth 1 seldepth 2 multipv 3 score mate -100 upperbound", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: Mate, Value: -100, Upperbound: true}}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 time 30", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, Time: 30}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 nodes 30", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, Nodes: 30}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 nps 30", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, Nps: 30}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 hashfull 30", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, HashFull: 30}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 tbhits 30", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, TbHits: 30}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 cpuload 30", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, CpuLoad: 30}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 string hello world", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, String: "hello world"}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 refutation e2e4 e7e5", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, Refutation: []string{"e2e4", "e7e5"}}},
	{"info depth 1 seldepth 2 multipv 3 score cp 100 currmove e2e4 currmovenumber 1", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, CurrentMove: "e2e4", CurrentMoveNumber: 1}},
	{"nfo depth 1 seldepth 2 multipv 3 score cp 100 currmove e2e4 currmovenumber 1 pv e2e4 e7e5", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, CurrentMove: "e2e4", CurrentMoveNumber: 1, Pv: []string{"e2e4", "e7e5"}}},
	{"nfo depth 1 seldepth 2 multipv 3 score cp 100 currmove e2e4 currmovenumber 1 pv e2e4 e7e5 refutation e2e4 e7e5", MoveInfo{Depth: 1, SelDepth: 2, MultiPv: 3, Score: Score{Type: CP, Value: 100}, CurrentMove: "e2e4", CurrentMoveNumber: 1, Pv: []string{"e2e4", "e7e5"}, Refutation: []string{"e2e4", "e7e5"}}},
}

func TestParseScore(t *testing.T) {
	for _, io := range scores {
		parts := strings.Split(io.in, " ")
		score := Score{}

		if n := parseScore(parts, 0, &score); n != io.tokens {
			t.Errorf("Expected %d, got %d", len(parts), n)
		}

		if !reflect.DeepEqual(score, io.out) {
			t.Errorf("Expected %v, got %v", io.out, score)
		}
	}
}

func TestParseInfoString(t *testing.T) {
	for _, io := range infostrs {
		info := parseInfoStr(io.in)

		if !reflect.DeepEqual(*info, io.out) {
			t.Errorf("Expected %v, got %v", io.out, *info)
		}
	}
}
