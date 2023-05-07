package run

import (
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/uci"
)

type model struct {
	gameData *data
	update   chan bool
	game     string
	color    string
}

type data struct {
	engine *uci.UCI
	client *com.Client
	moves  []string
	player string
	winner string
	reason string
	wtime  int
	btime  int
	err    error
	ttm    int
}

func (m model) isGameOver() bool {
	return m.gameData.winner != ""
}

func initModel(g string, c string, e *mgmt.EngineInstance) *model {
	engine, err := uci.New(e)

	if err != nil {
		return &model{gameData: &data{err: err}}
	}

	updates := make(chan bool)
	client, err := com.Connect()

	return &model{
		gameData: &data{engine: engine, client: client, err: err},
		update:   updates,
		game:     g,
		color:    c,
	}
}
