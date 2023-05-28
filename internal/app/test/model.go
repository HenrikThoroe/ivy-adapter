package test

import (
	"time"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com/testflow"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/schollz/progressbar/v3"
)

type state int

type searchMode int

const (
	connect state = iota
	wait
	play
	quit
)

const (
	searchTime searchMode = iota
	searchDepth
)

type model struct {
	bar     *progressbar.ProgressBar
	service *testService
	data    *data
	uptime  stopwatch.Model
}

type search struct {
	mode  searchMode
	value int
}

type options struct {
	hash    int
	threads int
}

type data struct {
	state       state
	err         error
	played      int
	session     string
	engines     [2]mgmt.EngineInstance
	search      [2]search
	options     [2]options
	concurrency int
}

func initModel() *model {
	client, err := com.Connect(conf.GetTestServerConfig().GetURL(), testflow.NewFlow())
	service := &testService{client: client}

	return &model{
		service: service,
		uptime:  stopwatch.NewWithInterval(time.Second),
		bar: progressbar.NewOptions(
			-1,
			progressbar.OptionSetDescription(""),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionSetElapsedTime(false),
			progressbar.OptionSetWidth(30),
		),
		data: &data{
			err:         err,
			state:       connect,
			concurrency: 1,
		},
	}
}
