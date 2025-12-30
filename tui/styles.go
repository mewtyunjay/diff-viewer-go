package tui

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	addColor    = lipgloss.Color("#2ECC71")
	deleteColor = lipgloss.Color("#E74C3C")
	addBg       = lipgloss.Color("#1E3A2F")
	deleteBg    = lipgloss.Color("#3A1E1E")
	contextFg   = lipgloss.Color("#AAAAAA")
	lineNumFg   = lipgloss.Color("#666666")
)

var (
	PanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtle)

	FocusedPanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(highlight)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(highlight).
			Padding(0, 1)

	TitleInactiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(subtle).
				Padding(0, 1)
)

var (
	AddLineStyle = lipgloss.NewStyle().
			Foreground(addColor).
			Background(addBg)

	DeleteLineStyle = lipgloss.NewStyle().
			Foreground(deleteColor).
			Background(deleteBg)

	ContextLineStyle = lipgloss.NewStyle().
				Foreground(contextFg)

	LineNumStyle = lipgloss.NewStyle().
			Foreground(lineNumFg).
			Width(5).
			Align(lipgloss.Right)

	EmptyLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333"))

	PlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#444444")).
				Background(lipgloss.Color("#1a1a1a"))
)

var (
	FileItemStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	FileItemSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Background(highlight).
				Foreground(lipgloss.Color("#FFFFFF"))

	AddCountStyle = lipgloss.NewStyle().
			Foreground(addColor)

	DelCountStyle = lipgloss.NewStyle().
			Foreground(deleteColor)
)

var HelpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#626262"))
