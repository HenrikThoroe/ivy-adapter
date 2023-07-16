package test

import (
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/com/testflow"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(m.service.register, m.uptime.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.data.err != nil {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case error:
		m.data.err = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.data.state = quit
			return m, tea.Quit
		default:
			return m, nil
		}
	case registerMsg:
		m.data.state = wait
		return m, m.service.awaitGameStart
	case startMsg:
		m.data.state = play
		m.data.session = msg.session
		m.data.engines = msg.engines
		m.data.search = msg.search
		m.data.options = msg.options
		m.data.concurrency = m.service.getConcurrency(msg.options)
		return m, func() tea.Msg {
			return m.service.dispatchGames(msg.batch, m.data)
		}
	case gameMsg:
		m.data.state = wait
		m.data.played += msg.gameCount
		m.service.client.Commands <- testflow.BuildReportCmd(m.data.session, msg.moves, msg.logs)
		return m, m.service.awaitGameStart
	}

	var cmd tea.Cmd
	m.uptime, cmd = m.uptime.Update(msg)
	return m, cmd
}
