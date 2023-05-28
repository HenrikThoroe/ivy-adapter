package test

import tea "github.com/charmbracelet/bubbletea"

func BuildTestViewModel() tea.Model {
	return *initModel()
}
