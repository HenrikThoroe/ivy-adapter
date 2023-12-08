package run

import (
	"github.com/HenrikThoroe/ivy-adapter/internal/app/instl"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/uci"
	tea "github.com/charmbracelet/bubbletea"
)

// SetupEngineIfc creates a new UCI interface based on the given path or
// installation name and version.
// If the path is empty, the installation view model will be shown.
// Otherwise the engine will be started with the given path.
// The function returns the UCI interface, the stdin and stdout channels and
// an error if one occurred.
func SetupEngineIfc(path string, name string, version string) (*uci.UCI, chan string, chan string, error) {
	var ifc *uci.UCI
	stdin := make(chan string)
	stdout := make(chan string)
	pipe := func(c chan string) func(s string) {
		return func(s string) {
			c <- s
		}
	}

	if path != "" {
		e, err := uci.NewFromExe(path, pipe(stdin), pipe(stdout))

		if err == nil {
			ifc = e
		} else {
			return nil, nil, nil, err
		}
	} else {
		model, inst := instl.BuildInstallationViewModel(name, version)

		if _, err := tea.NewProgram(model).Run(); err != nil {
			return nil, nil, nil, err
		}

		if e, err := uci.NewFromExe(inst.Path(), pipe(stdin), pipe(stdout)); err == nil {
			ifc = e
		} else {
			return nil, nil, nil, err
		}
	}

	return ifc, stdin, stdout, nil
}
