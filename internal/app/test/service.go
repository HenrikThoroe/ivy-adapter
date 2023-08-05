package test

import (
	"errors"
	"math"
	"runtime"
	"strconv"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com/testflow"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/sys"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/uci"
	tea "github.com/charmbracelet/bubbletea"
)

type registerMsg struct {
	id string
}

type startMsg struct {
	session string
	engines [2]mgmt.EngineInstance
	search  [2]search
	batch   int
	options [2]options
}

type gameMsg struct {
	gameCount int
	moves     []testflow.GameMoveHistory
	logs      []testflow.Log
}

type testService struct {
	client *com.Client
}

func (ts testService) register() tea.Msg {
	ts.client.Commands <- testflow.BuildRegisterCmd()

	select {
	case resp := <-ts.client.Messages:
		if rm, ok := resp.(testflow.RegisteredMsg); ok {
			return registerMsg{id: rm.Id}
		} else {
			return errors.New("did not receive registration confirmation")
		}
	case err := <-ts.client.Errors:
		return err
	}
}

func (ts testService) awaitGameStart() tea.Msg {
	select {
	case resp := <-ts.client.Messages:
		if sm, ok := resp.(testflow.StartMsg); ok {
			if len(sm.Suite.Engines) != 2 {
				return errors.New("did not receive two engines")
			}

			result := startMsg{
				session: sm.Session,
				engines: [2]mgmt.EngineInstance{},
				search:  [2]search{},
				batch:   sm.RecommendedBatchSize,
				options: [2]options{},
			}

			for idx, e := range sm.Suite.Engines {
				mode := searchTime
				engine, err := mgmt.BestMatch(e.Name, mgmt.Version{
					Major: e.Version.Major,
					Minor: e.Version.Minor,
					Patch: e.Version.Patch,
				})

				if err != nil {
					return err
				}

				if !engine.IsInstalled() {
					if err := mgmt.DownloadEngine(engine); err != nil {
						return err
					}
				}

				switch e.Time.Type {
				case "depth":
					mode = searchDepth
				case "movetime":
					mode = searchTime
				default:
					return errors.New("invalid search type")
				}

				result.engines[idx] = *engine
				result.search[idx] = search{
					mode:  mode,
					value: e.Time.Value,
				}
				result.options[idx] = options{
					hash:    e.Options.HashSize,
					threads: e.Options.Threads,
				}
			}

			return result
		} else {
			return errors.New("did not receive start confirmation")
		}
	case err := <-ts.client.Errors:
		return err
	}
}

func (ts testService) isGameOver(info *uci.MoveInfo) bool {
	noMove := info.Move == "(none)"
	mate := info.Score.Type == uci.Mate && info.Score.Value == 1

	return noMove || mate
}

func (ts testService) getConcurrency(options [2]options) int {
	device, _ := sys.DeviceInfo()
	cores := runtime.NumCPU()
	threads := int(math.Max(float64(options[0].threads), float64(options[1].threads)))
	requiredMemory := options[0].hash + options[1].hash + 512
	availableMemory := device.Memory / 1024 / 1024
	cpuLimit := cores / threads
	memLimit := availableMemory / requiredMemory
	limit := int(math.Min(float64(cpuLimit), float64(memLimit)))

	return limit
}

func (ts testService) dispatchGames(batch int, data *data) tea.Msg {
	cap := data.concurrency
	numGames := cap
	finished := 0
	active := 0
	gameChan := make(chan gameMsg)
	errChan := make(chan error)
	result := gameMsg{}

	if numGames < 1 {
		numGames = 1
	}

	for numGames < batch {
		numGames += cap
	}

	for i := 0; i < cap; i++ {
		active++
		go ts.awaitGame(gameChan, errChan, data)
	}

	for finished < numGames {
		select {
		case msg := <-gameChan:
			finished++
			active--
			result.gameCount += msg.gameCount
			result.moves = append(result.moves, msg.moves...)
			result.logs = append(result.logs, msg.logs...)
		case err := <-errChan:
			return err
		}

		if finished+active < numGames {
			active++
			go ts.awaitGame(gameChan, errChan, data)
		}
	}

	return result
}

func (ts testService) awaitGame(msg chan gameMsg, err chan error, data *data) {
	resp := ts.playGamePair(data)

	switch resp := resp.(type) {
	case error:
		err <- resp
	default:
		msg <- resp.(gameMsg)
	}
}

func (ts testService) playGamePair(data *data) tea.Msg {
	result := gameMsg{
		gameCount: 0,
		moves:     []testflow.GameMoveHistory{},
		logs:      []testflow.Log{},
	}

	resp1 := ts.playGame(data, false)

	switch resp1 := resp1.(type) {
	case error:
		return resp1
	case gameMsg:
		result.gameCount += resp1.gameCount
		result.moves = append(result.moves, resp1.moves...)
		result.logs = append(result.logs, resp1.logs...)
	}

	resp2 := ts.playGame(data, true)

	switch resp2 := resp2.(type) {
	case error:
		return resp2
	case gameMsg:
		result.gameCount += resp2.gameCount
		result.moves = append(result.moves, resp2.moves...)
		result.logs = append(result.logs, resp2.logs...)
	}

	return result
}

func (ts testService) playGame(data *data, swapColor bool) tea.Msg {
	ifc := [2]*uci.UCI{}
	maxMoves := 250
	moves := make([]string, 0, maxMoves)
	info := &uci.MoveInfo{}
	moveIdx := 0
	engineIdx := 0
	history := make([][]uci.MoveInfo, 2)
	logs := make([][]testflow.LogEntry, 2)

	for idx, e := range data.engines {
		logs[idx] = make([]testflow.LogEntry, 0, 1024)
		log := &logs[idx]

		scb := func(send string) {
			*log = append(*log, testflow.LogEntry{
				Type:  "send",
				Value: send,
			})
		}

		rcb := func(recv string) {
			*log = append(*log, testflow.LogEntry{
				Type:  "recv",
				Value: recv,
			})
		}

		u, err := uci.New(&e, scb, rcb)

		if err != nil {
			return err
		}

		ifc[idx] = u
	}

	if swapColor {
		engineIdx = 1
	}

	for idx, u := range ifc {
		u.Setup()

		if opt := u.GetOptionConfig("hash"); opt != nil {
			u.SetOption(opt.Response(strconv.Itoa(data.options[idx].hash)))
		}

		if opt := u.GetOptionConfig("threads"); opt != nil {
			u.SetOption(opt.Response(strconv.Itoa(data.options[idx].threads)))
		}

		u.Start()
	}

	for !ts.isGameOver(info) && moveIdx < maxMoves {
		ifc[engineIdx].SetMoves(moves...)
		info = ifc[engineIdx].GetMove(data.search[engineIdx].value)
		moves = append(moves, info.Move)
		history[engineIdx] = append(history[engineIdx], *info)
		engineIdx = (engineIdx + 1) % 2
		moveIdx++
	}

	for _, u := range ifc {
		u.Quit()
		u.Close()
	}

	return gameMsg{
		gameCount: 1,
		moves:     []testflow.GameMoveHistory{history},
		logs:      []testflow.Log{logs},
	}
}
