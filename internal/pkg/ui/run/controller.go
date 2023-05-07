package run

import (
	tea "github.com/charmbracelet/bubbletea"
)

func waitForUpdate(m model) tea.Cmd {
	return func() tea.Msg {
		return <-m.update
	}
}

func (m model) Init() tea.Cmd {
	if m.gameData.err != nil {
		return tea.Quit
	}

	go m.loop()

	return waitForUpdate(m)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case bool:
		if m.gameData.err != nil || m.isGameOver() {
			return m, tea.Quit
		}

		return m, waitForUpdate(m)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, cmd
}
