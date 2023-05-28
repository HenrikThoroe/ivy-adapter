package instl

import (
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/schollz/progressbar/v3"
)

// model is the data model for the installation view.
type model struct {
	engines        []mgmt.EngineInstance    // The available engines for installation
	selected       *mgmt.EngineInstance     // The selected engine for installation
	table          table.Model              // The table model for the available engines
	isLoading      bool                     // Whether the view is currently loading data
	bar            *progressbar.ProgressBar // The progress bar for the download
	err            error                    // Any error that occurred during the installation
	defaultEngine  string                   // The default engine name to install
	defaultVersion string                   // The default version name of the default engine to install
}

// initModel initializes the model with the default engine and version.
// If the default engine and version are not available, the user will be
// prompted to select an engine and version.
func initModel(defEng string, defVers string) model {
	var selected *mgmt.EngineInstance
	var e error

	loading := false
	bar := createProgressBar()

	if defEng != "" && defVers != "" {
		inst, err := mgmt.GetEngineInstance(defEng, defVers)

		if err == nil {
			selected = inst
			loading = true
			bar.Describe(string(downloadEngineDescription))
		} else {
			e = err
		}
	}

	return model{
		engines:        make([]mgmt.EngineInstance, 0),
		selected:       selected,
		table:          createTable(),
		isLoading:      loading,
		bar:            bar,
		err:            e,
		defaultEngine:  defEng,
		defaultVersion: defVers,
	}
}
