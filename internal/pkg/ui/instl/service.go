package instl

import (
	"time"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	tea "github.com/charmbracelet/bubbletea"
)

// fetchEngineConfig creates a command that fetches the available engines and
// returns an engineUpdateMsg with the result.
func fetchEngineConfig(engine string) tea.Cmd {
	return func() tea.Msg {
		engines, err := mgmt.GetAvailableEngines(engine)

		if err != nil {
			return errorMsg{err}
		}

		instances := make([]mgmt.EngineInstance, 0)

		for _, engine := range *engines {
			for _, inst := range engine.GetInstances() {
				instances = append(instances, inst)
			}
		}

		return engineUpdateMsg{instances}
	}
}

// downloadEngine creates a command that downloads the given engine and returns
// an engineDownloadedMsg with the result.
func downloadEngine(e *mgmt.EngineInstance) tea.Cmd {
	return func() tea.Msg {
		err := mgmt.DownloadEngine(e)

		if err != nil {
			return errorMsg{err}
		}

		return engineDownloadedMsg{e}
	}
}

// tick creates a command that sends a time.Time message every 10ms.
func tick() tea.Cmd {
	return tea.Every(time.Duration(1000*1000*10), func(t time.Time) tea.Msg {
		return time.Time(t)
	})
}
