package instl

import (
	"math"
	"strings"
	"time"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// engineUpdateMsg is a message that is sent when the available engines are updated.
type engineUpdateMsg struct {
	engines []engineInstance
}

// engineDownloadedMsg is a message that is sent when an engine has been downloaded.
type engineDownloadedMsg struct {
	engine *mgmt.EngineInstance
}

// errorMsg is a message that is sent when an error occurs.
type errorMsg struct {
	err error
}

// Init creates an initial command for the view model.
// If a default engine and version is set, it will try to download that engine.
// Otherwise it will fetch the available engines.
func (m model) Init() tea.Cmd {
	if m.selected != nil {
		return tea.Batch(tick(), downloadEngine(m.selected))
	}

	return tea.Batch(tick(), fetchEngineConfig(m.defaultEngine, m.defaultVersion))
}

// Update updates the view model.
// It handles messages from completed command and key presses.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case errorMsg:
		m.err = msg.err
		return m, tea.Quit
	case engineUpdateMsg:
		setAvailableEngines(&m, &msg)
		return m, nil
	case time.Time:
		if !m.isLoading || m.err != nil {
			return m, nil
		}

		m.bar.Add(1)
		return m, tick()
	case engineDownloadedMsg:
		m.isLoading = false
		m.downloadedEngine.Engine = msg.engine.Engine
		m.downloadedEngine.Version = msg.engine.Version
		m.downloadedEngine.Id = msg.engine.Id
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.table.Focused() {
				m.selected = &m.engines[m.table.Cursor()]
				m.isLoading = true
				m.bar.Describe(string(downloadEngineDescription))
				m.bar.Set(0)
				m.bar.Reset()
				return m, tea.Batch(tick(), downloadEngine(m.selected))
			}
		}
	}

	t, cmd := m.table.Update(msg)
	m.table = t

	return m, cmd
}

// setAvailableEngines sets the available engines in the model.
// It also sets the table rows and columns as well as the loading state.
func setAvailableEngines(m *model, msg *engineUpdateMsg) {
	rows := make([]table.Row, len(msg.engines))
	m.table.SetColumns([]table.Column{
		{Title: "Engine", Width: 15},
		{Title: "Version", Width: 10},
		{Title: "OS", Width: 10},
		{Title: "Arch", Width: 10},
		{Title: "Features", Width: 20},
	})

	for i, engine := range msg.engines {
		rows[i] = table.Row{
			engine.name,
			engine.version.String(mgmt.DotVersionStyle),
			engine.flavour.Os,
			engine.flavour.Arch,
			strings.Join(engine.flavour.Capabilities, ", "),
		}
	}

	m.engines = msg.engines
	m.table.SetRows(rows)
	m.table.SetHeight(int(math.Min(float64(len(rows)), 10.0)))
	m.isLoading = false
}
