package tui

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	// Diff colors
	addColor    = lipgloss.Color("#2ECC71")
	deleteColor = lipgloss.Color("#E74C3C")
	addBg       = lipgloss.Color("#1E3A2F")
	deleteBg    = lipgloss.Color("#3A1E1E")
	contextFg   = lipgloss.Color("#AAAAAA")
	lineNumFg   = lipgloss.Color("#666666")
)

// Panel styles
var (
	// Base panel style with rounded border
	PanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtle)

	// Focused panel has highlighted border
	FocusedPanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(highlight)

	// Panel title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(highlight).
			Padding(0, 1)

	// Unfocused title style
	TitleInactiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(subtle).
				Padding(0, 1)
)

// Diff line styles
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

	// Empty line style (for alignment)
	EmptyLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#333333"))
)

// File list styles
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

// Help style
var (
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)
