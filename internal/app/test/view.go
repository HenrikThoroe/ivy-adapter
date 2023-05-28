package test

import (
	"strconv"
	"time"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/mgmt"
	"github.com/charmbracelet/lipgloss"
)

type panelRow struct {
	label string
	value []string
}

type panel struct {
	title string
	rows  []panelRow
	width int
}

func (m model) View() string {
	if m.data.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true).
			Render("✘ Error: "+m.data.err.Error()) + "\n"
	}

	if m.data.state == quit {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true).
			Render("✘ Quit after "+strconv.Itoa(m.data.played)+" games") + "\n"
	}

	wrapper := lipgloss.NewStyle().
		Padding(1, 2)

	return wrapper.Render(m.createStatsPanel().build()) +
		"\n" +
		wrapper.Render(m.createEnginesPanel().build()) +
		"\n"
}

func (m model) createEnginesPanel() *panel {
	names := [2]string{
		"(none)",
		"(none)",
	}

	versions := [2]string{
		"(none)",
		"(none)",
	}

	modes := [2]string{
		"(none)",
		"(none)",
	}

	values := [2]string{
		"(none)",
		"(none)",
	}

	hashSizes := [2]string{
		"(none)",
		"(none)",
	}

	threads := [2]string{
		"(none)",
		"(none)",
	}

	if m.data.state == play {
		for idx, engine := range m.data.engines {
			names[idx] = engine.Engine
			versions[idx] = engine.Version.String(mgmt.DotVersionStyle)
			hashSizes[idx] = strconv.Itoa(m.data.options[idx].hash) + " MB"
			threads[idx] = strconv.Itoa(m.data.options[idx].threads)

			switch m.data.search[idx].mode {
			case searchDepth:
				modes[idx] = "Depth"
				values[idx] = strconv.Itoa(m.data.search[idx].value)
			case searchTime:
				modes[idx] = "Time"
				values[idx] = (time.Duration(m.data.search[idx].value) * time.Millisecond).String()
			}
		}
	}

	return &panel{
		title: "Engines",
		width: 60,
		rows: []panelRow{
			{
				label: "Name",
				value: names[:],
			},
			{
				label: "Version",
				value: versions[:],
			},
			{
				label: "Mode",
				value: modes[:],
			},
			{
				label: "Value",
				value: values[:],
			},
			{
				label: "Hash Size",
				value: hashSizes[:],
			},
			{
				label: "Threads",
				value: threads[:],
			},
		},
	}
}

func (m model) createStatsPanel() *panel {
	stateMsg := ""
	playedMsg := strconv.Itoa(m.data.played)
	uptimeMsg := m.uptime.View()
	concurrencyMsg := strconv.Itoa(m.data.concurrency)

	switch m.data.state {
	case connect:
		stateMsg = "Connecting to game server..."
	case play:
		stateMsg = "Playing Game"
	case wait:
		stateMsg = "Waiting for game to start..."
	}

	return &panel{
		title: "Stats",
		width: 60,
		rows: []panelRow{
			{
				label: "State",
				value: []string{stateMsg},
			},
			{
				label: "Played",
				value: []string{playedMsg},
			},
			{
				label: "Uptime",
				value: []string{uptimeMsg},
			},
			{
				label: "Concurrent Games",
				value: []string{concurrencyMsg},
			},
		},
	}
}

func (p panel) build() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("241")).
		Padding(0, 1)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#B949FF")).
		Bold(true).
		Align(lipgloss.Center).
		PaddingBottom(1)

	labelStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(20).
		Bold(true).
		Padding(0, 1)

	valueStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Padding(0, 1)

	var rows []string

	rows = append(rows, titleStyle.Render(p.title))

	for _, row := range p.rows {
		width := (p.width - 10) / len(row.value)
		values := []string{
			labelStyle.Render(row.label),
		}

		for _, value := range row.value {
			values = append(
				values,
				valueStyle.Width(width).Render(value),
			)
		}

		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Left, values...))
	}

	return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, rows...))
}
