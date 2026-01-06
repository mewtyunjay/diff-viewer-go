package tui

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	addColor    = lipgloss.Color("#2ECC71")
	deleteColor = lipgloss.Color("#E74C3C")
	contextFg   = lipgloss.Color("#AAAAAA")

	addBg    = lipgloss.Color("#1E3A2F")
	deleteBg = lipgloss.Color("#3A1E1E")
	lineNumFg = lipgloss.Color("#666666")

	// Intense highlight colors for changed portions within a line
	addHighlightBg    = lipgloss.Color("#2a5a3a")
	deleteHighlightBg = lipgloss.Color("#5a2a2a")
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

	// Highlight styles for word-level diff (changed portions within a line)
	AddHighlightStyle = lipgloss.NewStyle().
				Foreground(addColor).
				Background(addHighlightBg).
				Bold(true)

	DeleteHighlightStyle = lipgloss.NewStyle().
				Foreground(deleteColor).
				Background(deleteHighlightBg).
				Bold(true)
)

var (
	FileItemStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	// Darker blue background for selection (like lazygit)
	FileItemSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Background(lipgloss.Color("#3d59a1")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)

	AddCountStyle = lipgloss.NewStyle().
			Foreground(addColor)

	DelCountStyle = lipgloss.NewStyle().
			Foreground(deleteColor)

	// Status indicator styles (for tree view)
	StatusModifiedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e5c07b")) // Yellow

	StatusNewStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#56b6c2")) // Cyan

	StatusDeletedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e06c75")) // Red

	// Expand/collapse indicators for tree
	ExpandedIndicator  = "▼"
	CollapsedIndicator = "▶"
)

var HelpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#626262"))

// Modal styles for commit dialog
var (
	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(1, 2)

	ModalTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			MarginBottom(1)

	ModalHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			MarginTop(1)

	ModalErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E74C3C")).
			MarginTop(1)

	// Staged file indicator style
	StatusStagedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#2ECC71")). // Green for staged
				Bold(true)
)
