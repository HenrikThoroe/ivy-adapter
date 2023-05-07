package run

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) renderStats() string {
	lastMove := "(none)"

	if len(m.gameData.moves) > 0 {
		lastMove = m.gameData.moves[len(m.gameData.moves)-1]
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("57")).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Width(18).
		MarginTop(1).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		MarginTop(1).
		Bold(false)

	wduration := time.Duration(m.gameData.wtime) * time.Millisecond
	bduration := time.Duration(m.gameData.btime) * time.Millisecond

	if wduration > 10*time.Second {
		wduration = wduration.Round(time.Second)
	}

	if bduration > 10*time.Second {
		bduration = bduration.Round(time.Second)
	}

	if wduration < 0 {
		wduration = 0
	}

	if bduration < 0 {
		bduration = 0
	}

	info := [][]string{
		{"Game ID:", m.game},
		{"Player ID:", m.gameData.player},
		{"Ping:", strconv.FormatInt(m.gameData.client.Ping(), 10) + "ms"},
		{"Color:", m.color},
		{"Time:", fmt.Sprintf("%s / %s (last move: %s)", wduration, bduration, time.Duration(m.gameData.ttm)*time.Millisecond)},
		{"Last Move:", lastMove},
		{"Winner:", m.gameData.winner + " (" + m.gameData.reason + ")"},
	}

	if m.gameData.winner == "" {
		info[6][1] = "(none)"
	}

	var rows []string

	for _, row := range info {
		rows = append(rows,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				labelStyle.Render(row[0]),
				valueStyle.Render(row[1]),
			),
		)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return lipgloss.JoinVertical(lipgloss.Center,
		headerStyle.Render("Stats"),
		content,
	)
}

func (m *model) renderBoard() string {
	posToIdx := func(pos string) (int, int) {
		file := int(pos[0] - 'a')
		rank := int(pos[1] - '1')
		row := 7 - rank
		col := file

		return row, col
	}

	board := [][]string{
		{"r", "n", "b", "q", "k", "b", "n", "r"},
		{"p", "p", "p", "p", "p", "p", "p", "p"},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{"P", "P", "P", "P", "P", "P", "P", "P"},
		{"R", "N", "B", "Q", "K", "B", "N", "R"},
	}

	for _, move := range m.gameData.moves {
		promotion := ""

		if len(move) == 5 {
			promotion = move[4:]
		}

		startX, startY := posToIdx(move[:2])
		endX, endY := posToIdx(move[2:])
		piece := board[startX][startY]

		if piece == "K" || piece == "k" {
			if endY-startY == 2 {
				board[startX][5] = board[startX][7]
				board[startX][7] = " "
			} else if endY-startY == -2 {
				board[startX][3] = board[startX][0]
				board[startX][0] = " "
			}
		}

		if promotion != "" {
			if piece == "P" {
				promotion = strings.ToUpper(promotion)
			} else {
				promotion = strings.ToLower(promotion)
			}

			board[endX][endY] = promotion
		} else {
			board[endX][endY] = board[startX][startY]
		}

		board[startX][startY] = " "
	}

	res := ""

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("57")).
		MarginBottom(1).
		Bold(true)

	fieldStyle := lipgloss.NewStyle().
		Width(3).
		Height(1).
		Align(lipgloss.Center)

	darkStyle := fieldStyle.Copy().
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("15"))

	lightStyle := fieldStyle.Copy().
		Background(lipgloss.Color("15")).
		Foreground(lipgloss.Color("57"))

	boardIndexStyle := lipgloss.NewStyle().
		Width(3).
		Height(1).
		Align(lipgloss.Center).
		Bold(true)

	for i, row := range board {
		res += boardIndexStyle.Render(strconv.Itoa(8 - i))

		for j, piece := range row {
			if (i+j)%2 == 0 {
				res += lightStyle.Render(piece)
			} else {
				res += darkStyle.Render(piece)
			}
		}
		res += "\n"
	}

	res += boardIndexStyle.Render(" ")

	for i := 0; i < 8; i++ {
		res += boardIndexStyle.Render(string(rune('a' + i)))
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		headerStyle.Render("Board"),
		res,
	)
}

func (m model) View() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 4)

	if m.gameData.err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true).
			Render("✘ Error: "+m.gameData.err.Error()) + "\n"
	}

	return boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			m.renderStats(),
			"\n\n",
			m.renderBoard(),
		),
	) + "\n"
}
