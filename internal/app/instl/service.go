package instl

import (
	"time"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	tea "github.com/charmbracelet/bubbletea"
)

type engineInstance struct {
	name    string
	version mgmt.Version
	flavour mgmt.Flavour
}

// fetchEngineConfig creates a command that fetches the available engines and
// returns an engineUpdateMsg with the result.
func fetchEngineConfig(engine string, version string) tea.Cmd {
	return func() tea.Msg {
		engines, err := mgmt.GetAvailableEngines(engine)

		if err != nil {
			return errorMsg{err}
		}

		instances := make([]engineInstance, 0)
		vers, err := mgmt.ParseVersion(version, mgmt.DotVersionStyle)

		for _, engine := range *engines {
			for _, vari := range engine.Variations {
				if err == nil && (vers.Major != vari.Version.Major || vers.Minor != vari.Version.Minor || vers.Patch != vari.Version.Patch) {
					continue
				}

				for _, flav := range vari.Flavours {
					instances = append(instances, engineInstance{
						name:    engine.Name,
						version: vari.Version,
						flavour: flav,
					})
				}
			}
		}

		return engineUpdateMsg{
			engines: instances,
		}
	}
}

// downloadEngine creates a command that downloads the given engine and returns
// an engineDownloadedMsg with the result.
func downloadEngine(e *engineInstance) tea.Cmd {
	return func() tea.Msg {
		inst := mgmt.EngineInstance{
			Engine:  e.name,
			Version: e.version,
			Id:      e.flavour.Id,
		}

		err := mgmt.DownloadEngine(&inst)

		if err != nil {
			return errorMsg{err}
		}

		return engineDownloadedMsg{&inst}
	}
}

// tick creates a command that sends a time.Time message every 10ms.
func tick() tea.Cmd {
	return tea.Every(time.Duration(1000*1000*10), func(t time.Time) tea.Msg {
		return time.Time(t)
	})
}
