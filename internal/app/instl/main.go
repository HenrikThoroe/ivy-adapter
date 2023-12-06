// Package instl provides the installation view.
// It is used by the install command to display the installation view.
// The installation view is used to select an engine and version to install.
// The installation view is also used to display the progress of the installation.
package instl

import (
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	tea "github.com/charmbracelet/bubbletea"
)

// BuildInstallationViewModel returns a new tea.Model for the installation view.
// It takes the engine and version as arguments. If the engine is not specified,
// the user will be prompted to select one. If the version is not specified, the
// user will be prompted to select one.
func BuildInstallationViewModel(engine string, version string) (tea.Model, *mgmt.EngineInstance) {
	model := initModel(engine, version)

	return model, model.downloadedEngine
}
