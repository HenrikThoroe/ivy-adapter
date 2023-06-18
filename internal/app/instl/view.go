package instl

import (
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/schollz/progressbar/v3"
)

// progressBarDescription is a type alias for the description of the progress bar.
type progressBarDescription string

const (
	// fetchEngineDescription is the description of the progress bar when fetching the engine list.
	fetchEngineDescription progressBarDescription = "Fetching engine list..."
	// downloadEngineDescription is the description of the progress bar when downloading the engine.
	downloadEngineDescription progressBarDescription = "Downloading engine..."
)

// baseStyle is the base style of the table.
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// headerStyle is the style of the table header.
var headerStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	BorderBottom(true).
	Bold(true)

// selectedStyle is the style of the selected table row.
var selectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("229")).
	Background(lipgloss.Color("57")).
	Bold(true)

// errorStyle is the style of the error message.
var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("196")).
	Bold(true)

// promptStyle is the style of the prompt.
var promptStyle = lipgloss.NewStyle().Bold(true)

var checkMarkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))

var crossStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

// View composes the TUI view based on the current state of the view model.
func (m model) View() string {
	if !m.isLoading {
		clearProgressBar(m.bar)
	}

	if m.err != nil {
		clearProgressBar(m.bar)
		return crossStyle.Render("✘ ") + errorStyle.Render(m.err.Error()) + "\n"
	}

	if m.isLoading {
		return m.bar.String()
	}

	if m.selected != nil {
		clearProgressBar(m.bar)
		prompt := "Downloaded " + m.selected.name + " @ " + m.selected.version.String(mgmt.DotVersionStyle)

		return checkMarkStyle.Render("✔ ") +
			promptStyle.Render(prompt) +
			"\n"
	}

	return promptStyle.Render("Please select an engine to install:") +
		"\n" +
		baseStyle.Render(applyTableStyles(&m.table).View()) +
		"\n"
}

// clearProgressBar removes the progress bar from the screen.
func clearProgressBar(bar *progressbar.ProgressBar) {
	bar.Clear()
	bar.Exit()
}

// applyTableStyles applies styles to the table.
func applyTableStyles(t *table.Model) *table.Model {
	tableStyle := table.DefaultStyles()
	tableStyle.Header = headerStyle.Padding(0, 1)
	tableStyle.Selected = selectedStyle
	tableStyle.Cell = lipgloss.NewStyle().Padding(0, 1)
	t.SetStyles(tableStyle)
	return t
}

// createProgressBar creates a new progress bar.
// By default the progress bar has the description for fetching the engine list.
// The progress bar acts as a spinner.
func createProgressBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription(string(fetchEngineDescription)),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetWidth(30),
	)
}

// createTable creates a new table without any data.
// Set rows and columns before drawing the table.
func createTable() table.Model {
	return table.New(
		table.WithHeight(10),
		table.WithFocused(true),
	)
}
