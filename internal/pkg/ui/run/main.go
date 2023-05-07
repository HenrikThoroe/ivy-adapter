package run

import (
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	tea "github.com/charmbracelet/bubbletea"
)

func BuildRunViewModel(g string, c string, e *mgmt.EngineInstance) tea.Model {
	return *initModel(g, c, e)
}
