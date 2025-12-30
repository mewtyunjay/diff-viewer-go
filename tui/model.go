package tui

import (
	"fmt"
	"strings"

	"diff-tui/diff"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FocusedPanel tracks which panel currently has focus
type FocusedPanel int

const (
	FocusFileList FocusedPanel = iota
	FocusLeftDiff
	FocusRightDiff
)

// Model is the main application model
type Model struct {
	width  int
	height int

	focused FocusedPanel
	keys    KeyMap

	// File list state
	files       []diff.FileDiff
	selectedIdx int

	// Diff viewports
	leftViewport  viewport.Model
	rightViewport viewport.Model

	// Synchronized scroll
	syncScroll bool

	ready bool
}

// NewModel creates a new TUI model with the given files
func NewModel(files []diff.FileDiff) Model {
	return Model{
		files:      files,
		keys:       DefaultKeyMap,
		syncScroll: true,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Tab):
			m.focused = (m.focused + 1) % 3

		case key.Matches(msg, m.keys.ShiftTab):
			m.focused = (m.focused + 2) % 3

		case key.Matches(msg, m.keys.SyncToggle):
			m.syncScroll = !m.syncScroll

		case key.Matches(msg, m.keys.Up):
			if m.focused == FocusFileList {
				if m.selectedIdx > 0 {
					m.selectedIdx--
					m.updateDiffContent()
				}
			} else {
				m.scrollUp(1)
			}

		case key.Matches(msg, m.keys.Down):
			if m.focused == FocusFileList {
				if m.selectedIdx < len(m.files)-1 {
					m.selectedIdx++
					m.updateDiffContent()
				}
			} else {
				m.scrollDown(1)
			}

		case key.Matches(msg, m.keys.PageUp):
			if m.focused != FocusFileList {
				m.scrollUp(m.leftViewport.Height)
			}

		case key.Matches(msg, m.keys.PageDown):
			if m.focused != FocusFileList {
				m.scrollDown(m.leftViewport.Height)
			}

		case key.Matches(msg, m.keys.HalfPageUp):
			if m.focused != FocusFileList {
				m.scrollUp(m.leftViewport.Height / 2)
			}

		case key.Matches(msg, m.keys.HalfPageDown):
			if m.focused != FocusFileList {
				m.scrollDown(m.leftViewport.Height / 2)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateViewportSizes()
		m.updateDiffContent()
		m.ready = true
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) scrollUp(lines int) {
	if m.focused == FocusLeftDiff || m.syncScroll {
		m.leftViewport.SetYOffset(max(0, m.leftViewport.YOffset-lines))
	}
	if m.focused == FocusRightDiff || m.syncScroll {
		m.rightViewport.SetYOffset(max(0, m.rightViewport.YOffset-lines))
	}
}

func (m *Model) scrollDown(lines int) {
	if m.focused == FocusLeftDiff || m.syncScroll {
		m.leftViewport.SetYOffset(min(m.leftViewport.TotalLineCount()-m.leftViewport.Height, m.leftViewport.YOffset+lines))
	}
	if m.focused == FocusRightDiff || m.syncScroll {
		m.rightViewport.SetYOffset(min(m.rightViewport.TotalLineCount()-m.rightViewport.Height, m.rightViewport.YOffset+lines))
	}
}

func (m *Model) updateViewportSizes() {
	// Calculate panel widths
	// Total width minus borders (2 chars per panel = 6 total)
	availableWidth := m.width - 6
	fileListWidth := availableWidth * 20 / 100
	diffPanelWidth := (availableWidth - fileListWidth) / 2

	// Height minus borders and title
	panelHeight := m.height - 4

	m.leftViewport = viewport.New(diffPanelWidth-2, panelHeight)
	m.rightViewport = viewport.New(diffPanelWidth-2, panelHeight)
}

func (m *Model) updateDiffContent() {
	if len(m.files) == 0 {
		return
	}

	file := m.files[m.selectedIdx]

	// Calculate available width for content
	availableWidth := m.width - 6
	fileListWidth := availableWidth * 20 / 100
	diffPanelWidth := (availableWidth - fileListWidth) / 2
	contentWidth := diffPanelWidth - 8 // Account for line numbers and padding

	m.leftViewport.SetContent(m.renderDiffLines(file.LeftLines, contentWidth, true))
	m.rightViewport.SetContent(m.renderDiffLines(file.RightLines, contentWidth, false))
}

func (m *Model) renderDiffLines(lines []diff.Line, width int, isLeft bool) string {
	var sb strings.Builder

	for i, line := range lines {
		lineNum := fmt.Sprintf("%4d ", i+1)

		// Determine style based on line type
		var style lipgloss.Style
		var numStyle lipgloss.Style
		content := line.Content

		switch line.Type {
		case diff.Add:
			style = AddLineStyle
			numStyle = LineNumStyle.Foreground(lipgloss.Color("#2ECC71"))
		case diff.Delete:
			style = DeleteLineStyle
			numStyle = LineNumStyle.Foreground(lipgloss.Color("#E74C3C"))
		default:
			style = ContextLineStyle
			numStyle = LineNumStyle
		}

		// For placeholder lines (filler lines in side-by-side view)
		if line.Type == diff.Placeholder {
			numStyle = LineNumStyle.Foreground(lipgloss.Color("#333333"))
			lineNum = "     "
			// Use filler pattern instead of blank space
			content = strings.Repeat("░", width)
			style = PlaceholderStyle
		} else if len(content) > width {
			// Truncate content to fit width
			content = content[:width-1] + "~"
		} else {
			// Pad content to fit width
			content = content + strings.Repeat(" ", width-len(content))
		}

		sb.WriteString(numStyle.Render(lineNum))
		sb.WriteString(style.Render(content))
		sb.WriteString("\n")
	}

	return sb.String()
}

// View implements tea.Model
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Calculate panel widths
	availableWidth := m.width - 6
	fileListWidth := availableWidth * 20 / 100
	diffPanelWidth := (availableWidth - fileListWidth) / 2

	// Panel height
	panelHeight := m.height - 2

	// Render panels
	leftPanel := m.renderFileListPanel(fileListWidth, panelHeight)
	middlePanel := m.renderDiffPanel("Original", m.leftViewport.View(), diffPanelWidth, panelHeight, m.focused == FocusLeftDiff)
	rightPanel := m.renderDiffPanel("Modified", m.rightViewport.View(), diffPanelWidth, panelHeight, m.focused == FocusRightDiff)

	// Join panels horizontally
	main := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, middlePanel, rightPanel)

	return main
}

func (m Model) renderFileListPanel(width, height int) string {
	isFocused := m.focused == FocusFileList

	// Title
	var title string
	if isFocused {
		title = TitleStyle.Render("Files")
	} else {
		title = TitleInactiveStyle.Render("Files")
	}
	title = lipgloss.PlaceHorizontal(width-2, lipgloss.Left, title)

	// File list content
	var content strings.Builder
	for i, file := range m.files {
		var style lipgloss.Style
		if i == m.selectedIdx {
			style = FileItemSelectedStyle
		} else {
			style = FileItemStyle
		}

		// File name
		name := file.Name
		if len(name) > width-12 {
			name = name[:width-15] + "..."
		}

		// Add/delete counts
		counts := fmt.Sprintf(" %s %s",
			AddCountStyle.Render(fmt.Sprintf("+%d", file.AddCount)),
			DelCountStyle.Render(fmt.Sprintf("-%d", file.DelCount)),
		)

		// Indicator for new/deleted files
		indicator := ""
		if file.IsNew {
			indicator = AddCountStyle.Render("[new] ")
		} else if file.IsDeleted {
			indicator = DelCountStyle.Render("[del] ")
		}

		line := fmt.Sprintf("%s%s%s", indicator, name, counts)
		line = style.Width(width - 2).Render(line)
		content.WriteString(line + "\n")
	}

	// Pad remaining height
	contentLines := len(m.files)
	remainingHeight := height - 3 - contentLines
	if remainingHeight > 0 {
		content.WriteString(strings.Repeat("\n", remainingHeight))
	}

	// Status line
	syncStatus := "sync: on"
	if !m.syncScroll {
		syncStatus = "sync: off"
	}
	statusLine := HelpStyle.Render(fmt.Sprintf(" %s | q: quit", syncStatus))

	// Build panel content
	panelContent := lipgloss.JoinVertical(lipgloss.Left,
		title,
		content.String(),
		statusLine,
	)

	// Apply panel style
	var panelStyle lipgloss.Style
	if isFocused {
		panelStyle = FocusedPanelStyle
	} else {
		panelStyle = PanelStyle
	}

	return panelStyle.
		Width(width).
		Height(height).
		Render(panelContent)
}

func (m Model) renderDiffPanel(title string, content string, width, height int, isFocused bool) string {
	// Title
	var titleRendered string
	if isFocused {
		titleRendered = TitleStyle.Render(title)
	} else {
		titleRendered = TitleInactiveStyle.Render(title)
	}
	titleRendered = lipgloss.PlaceHorizontal(width-2, lipgloss.Left, titleRendered)

	// Scroll indicator
	var scrollInfo string
	if isFocused || m.syncScroll {
		vp := m.leftViewport
		if title == "Modified" {
			vp = m.rightViewport
		}
		scrollPct := 0
		if vp.TotalLineCount() > 0 {
			scrollPct = int(float64(vp.YOffset) / float64(max(1, vp.TotalLineCount()-vp.Height)) * 100)
		}
		scrollInfo = HelpStyle.Render(fmt.Sprintf(" %d%% ", min(100, scrollPct)))
	}

	// Build panel content
	panelContent := lipgloss.JoinVertical(lipgloss.Left,
		titleRendered,
		content,
		scrollInfo,
	)

	// Apply panel style
	var panelStyle lipgloss.Style
	if isFocused {
		panelStyle = FocusedPanelStyle
	} else {
		panelStyle = PanelStyle
	}

	return panelStyle.
		Width(width).
		Height(height).
		Render(panelContent)
}
