package tui

import "github.com/charmbracelet/lipgloss"

var (
	colRed  = lipgloss.Color("#ff5a52")
	colTeal = lipgloss.Color("#5ffbe0")
	colBg   = lipgloss.Color("#0b0203")
)

const contentWidth = 56

var (
	frameStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colTeal).
			BorderBackground(colBg).
			Background(colBg).
			Padding(1, 2)

	headerStyle = lipgloss.NewStyle().Foreground(colTeal).Background(colBg).Bold(true)
	ruleStyle   = lipgloss.NewStyle().Foreground(colTeal).Background(colBg)
	promptStyle = lipgloss.NewStyle().Foreground(colRed).Background(colBg)
	rowStyle    = lipgloss.NewStyle().Foreground(colRed).Background(colBg)
	rowSelStyle = lipgloss.NewStyle().Foreground(colBg).Background(colRed).Bold(true)
	dimStyle    = lipgloss.NewStyle().Foreground(colRed).Background(colBg).Faint(true)
	tealStyle   = lipgloss.NewStyle().Foreground(colTeal).Background(colBg)
	okStyle     = lipgloss.NewStyle().Foreground(colTeal).Background(colBg).Bold(true)
	footerStyle = lipgloss.NewStyle().Foreground(colTeal).Background(colBg).Faint(true)

	bgFill = lipgloss.NewStyle().Background(colBg)
)
