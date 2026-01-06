package tui

import (
	"context"
	"fmt"
	"strings"

	"diff-tui/diff"
	"diff-tui/parser"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FocusedPanel int

const (
	FocusFileList FocusedPanel = iota
	FocusLeftDiff
	FocusRightDiff
)

// main application model
type Model struct {
	width  int
	height int

	focused FocusedPanel
	keys    KeyMap

	files        []diff.FileDiff
	treeRoots    []*TreeNode   // Root nodes of file tree
	visibleNodes []*TreeNode   // Flattened visible nodes for navigation
	selectedIdx  int           // Index in visibleNodes
	rootName     string        // Name of the root folder (displayed in tree)

	leftViewport  viewport.Model
	rightViewport viewport.Model

	syncScroll bool

	ready bool

	// Git operations
	gitRunner   *parser.GitRunner
	diffArgs    []string        // Original diff args for refresh
	stagedFiles map[string]bool // Track staged status by filepath

	// Commit modal state
	commitModalActive bool
	commitInput       textinput.Model
	commitError       string
}

// creates a new TUI model with the given files
func NewModel(files []diff.FileDiff, gitRunner *parser.GitRunner, diffArgs []string, rootName string) Model {
	treeRoots := BuildTree(files, rootName)
	visibleNodes := FlattenVisible(treeRoots)

	// Find initial selection (first file, not directory)
	selectedIdx := 0
	for i, node := range visibleNodes {
		if node.IsFile() {
			selectedIdx = i
			break
		}
	}

	// Initialize commit input
	ti := textinput.New()
	ti.Placeholder = "Enter commit message..."
	ti.CharLimit = 200
	ti.Width = 50

	// Initialize staged files map
	stagedFiles := make(map[string]bool)

	// Load currently staged files
	if gitRunner != nil {
		ctx := context.Background()
		staged, err := gitRunner.GetStagedFiles(ctx)
		if err == nil {
			for _, f := range staged {
				stagedFiles[f] = true
			}
		}
	}

	return Model{
		files:        files,
		treeRoots:    treeRoots,
		visibleNodes: visibleNodes,
		selectedIdx:  selectedIdx,
		rootName:     rootName,
		keys:         DefaultKeyMap,
		syncScroll:   true,
		gitRunner:    gitRunner,
		diffArgs:     diffArgs,
		stagedFiles:  stagedFiles,
		commitInput:  ti,
	}
}

// implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle commit modal input first
	if m.commitModalActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// Execute commit
				if m.commitInput.Value() != "" {
					m.executeCommit()
				}
				return m, nil
			case "esc":
				// Close modal
				m.closeCommitModal()
				return m, nil
			default:
				// Pass to text input
				var cmd tea.Cmd
				m.commitInput, cmd = m.commitInput.Update(msg)
				return m, cmd
			}
		}
		return m, nil
	}

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

		case key.Matches(msg, m.keys.Stage):
			if m.focused == FocusFileList {
				m.toggleStaging()
			}

		case key.Matches(msg, m.keys.Commit):
			if m.focused == FocusFileList && m.hasStagedFiles() {
				m.openCommitModal()
			}

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
				if m.selectedIdx < len(m.visibleNodes)-1 {
					m.selectedIdx++
					m.updateDiffContent()
				}
			} else {
				m.scrollDown(1)
			}

		case key.Matches(msg, m.keys.Left):
			if m.focused == FocusFileList {
				m.handleTreeLeft()
			}

		case key.Matches(msg, m.keys.Right):
			if m.focused == FocusFileList {
				m.handleTreeRight()
			}

		case key.Matches(msg, m.keys.Enter):
			if m.focused == FocusFileList {
				m.handleTreeEnter()
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

// handleTreeLeft collapses directory or goes to parent
func (m *Model) handleTreeLeft() {
	if len(m.visibleNodes) == 0 || m.selectedIdx >= len(m.visibleNodes) {
		return
	}

	node := m.visibleNodes[m.selectedIdx]

	if node.IsDirectory() && node.Expanded {
		// Collapse directory
		node.ToggleExpanded()
		m.refreshVisibleNodes()
	} else if node.Parent != nil {
		// Go to parent
		m.selectNode(node.Parent)
	}
}

// handleTreeRight expands directory
func (m *Model) handleTreeRight() {
	if len(m.visibleNodes) == 0 || m.selectedIdx >= len(m.visibleNodes) {
		return
	}

	node := m.visibleNodes[m.selectedIdx]

	if node.IsDirectory() && !node.Expanded {
		node.ToggleExpanded()
		m.refreshVisibleNodes()
	}
}

// handleTreeEnter toggles directory or selects file
func (m *Model) handleTreeEnter() {
	if len(m.visibleNodes) == 0 || m.selectedIdx >= len(m.visibleNodes) {
		return
	}

	node := m.visibleNodes[m.selectedIdx]

	if node.IsDirectory() {
		node.ToggleExpanded()
		m.refreshVisibleNodes()
	}
	// For files, Enter could switch focus to diff panel (optional enhancement)
}

// refreshVisibleNodes rebuilds the visible nodes list after expand/collapse
func (m *Model) refreshVisibleNodes() {
	currentNode := m.visibleNodes[m.selectedIdx]
	m.visibleNodes = FlattenVisible(m.treeRoots)

	// Try to keep the same node selected
	for i, node := range m.visibleNodes {
		if node == currentNode {
			m.selectedIdx = i
			return
		}
	}

	// If node is no longer visible (collapsed parent), select first visible
	if m.selectedIdx >= len(m.visibleNodes) {
		m.selectedIdx = max(0, len(m.visibleNodes)-1)
	}
}

// selectNode finds and selects a specific node
func (m *Model) selectNode(target *TreeNode) {
	for i, node := range m.visibleNodes {
		if node == target {
			m.selectedIdx = i
			m.updateDiffContent()
			return
		}
	}
}

func (m *Model) updateViewportSizes() {
	// calculate panel widths
	// total width minus borders (2 chars per panel = 6 total)
	availableWidth := m.width - 6
	fileListWidth := availableWidth * 20 / 100
	diffPanelWidth := (availableWidth - fileListWidth) / 2

	// height minus borders and title
	panelHeight := m.height - 4

	m.leftViewport = viewport.New(diffPanelWidth-2, panelHeight)
	m.rightViewport = viewport.New(diffPanelWidth-2, panelHeight)
}

func (m *Model) updateDiffContent() {
	if len(m.visibleNodes) == 0 || m.selectedIdx >= len(m.visibleNodes) {
		return
	}

	node := m.visibleNodes[m.selectedIdx]

	// Only show diff for files, not directories
	if node.File == nil {
		m.leftViewport.SetContent("")
		m.rightViewport.SetContent("")
		return
	}

	file := node.File

	// calculate available width for content
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

		// Determine styles based on line type
		var baseStyle, highlightStyle lipgloss.Style
		var numStyle lipgloss.Style

		switch line.Type {
		case diff.Add:
			baseStyle = AddLineStyle
			highlightStyle = AddHighlightStyle
			numStyle = LineNumStyle.Foreground(lipgloss.Color("#2ECC71"))
		case diff.Delete:
			baseStyle = DeleteLineStyle
			highlightStyle = DeleteHighlightStyle
			numStyle = LineNumStyle.Foreground(lipgloss.Color("#E74C3C"))
		default:
			baseStyle = ContextLineStyle
			highlightStyle = ContextLineStyle // No highlight for context
			numStyle = LineNumStyle
		}

		// Handle placeholder lines (filler lines in side-by-side view)
		if line.Type == diff.Placeholder {
			numStyle = LineNumStyle.Foreground(lipgloss.Color("#333333"))
			lineNum = "     "
			content := strings.Repeat("░", width)
			sb.WriteString(numStyle.Render(lineNum))
			sb.WriteString(PlaceholderStyle.Render(content))
			sb.WriteString("\n")
			continue
		}

		// Render line number
		sb.WriteString(numStyle.Render(lineNum))

		// Render content - with segments if available (word-level diff)
		if len(line.Segments) > 0 {
			sb.WriteString(m.renderSegmentedLine(line.Segments, baseStyle, highlightStyle, width))
		} else {
			// No segments - render entire content with base style
			content := line.Content
			if len(content) > width {
				content = content[:width-1] + "~"
			} else {
				content = content + strings.Repeat(" ", width-len(content))
			}
			sb.WriteString(baseStyle.Render(content))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderSegmentedLine renders a line with mixed highlighting for word-level diff
func (m *Model) renderSegmentedLine(segments []diff.Segment, baseStyle, highlightStyle lipgloss.Style, width int) string {
	var sb strings.Builder
	totalLen := 0

	for _, seg := range segments {
		if totalLen >= width {
			break
		}

		text := seg.Text
		remaining := width - totalLen

		// Truncate if needed
		if len(text) > remaining {
			text = text[:remaining-1] + "~"
		}

		if seg.Highlighted {
			sb.WriteString(highlightStyle.Render(text))
		} else {
			sb.WriteString(baseStyle.Render(text))
		}
		totalLen += len(text)
	}

	// Pad remaining width with base style
	if totalLen < width {
		sb.WriteString(baseStyle.Render(strings.Repeat(" ", width-totalLen)))
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

	// Overlay commit modal if active
	if m.commitModalActive {
		return m.renderCommitModal(main)
	}

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

	// File list content - render tree
	var content strings.Builder
	contentWidth := width - 4 // Account for borders and padding

	for i, node := range m.visibleNodes {
		isSelected := i == m.selectedIdx
		line := m.renderTreeNode(node, contentWidth, isSelected)
		content.WriteString(line + "\n")
	}

	// Pad remaining height
	contentLines := len(m.visibleNodes)
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

// renderTreeNode renders a single tree node with proper indentation
func (m Model) renderTreeNode(node *TreeNode, width int, isSelected bool) string {
	var sb strings.Builder

	// Indentation (2 spaces per depth level)
	indent := strings.Repeat("  ", node.Depth)
	sb.WriteString(indent)

	if node.IsDirectory() {
		// Directory: show expand/collapse indicator and name
		if node.Expanded {
			sb.WriteString(ExpandedIndicator + " ")
		} else {
			sb.WriteString(CollapsedIndicator + " ")
		}
		sb.WriteString(node.Name + "/")
	} else {
		// File: show status, name, and counts
		status := m.getFileStatus(node, isSelected)
		sb.WriteString(status + " ")
		sb.WriteString(node.Name)

		if node.File != nil {
			counts := fmt.Sprintf(" +%d -%d", node.File.AddCount, node.File.DelCount)
			sb.WriteString(counts)
		}
	}

	lineContent := sb.String()

	if isSelected {
		return FileItemSelectedStyle.Width(width).Render(lineContent)
	}
	return FileItemStyle.Width(width).Render(lineContent)
}

// getFileStatus returns the status indicator for a file
func (m Model) getFileStatus(node *TreeNode, isSelected bool) string {
	if node.File == nil {
		return "  "
	}

	// Check if file is staged
	isStaged := m.stagedFiles[node.File.Name]

	// When selected, return plain text (no ANSI codes) so selection style applies uniformly
	if isSelected {
		if isStaged {
			return "S "
		}
		if node.File.IsNew {
			return "??"
		} else if node.File.IsDeleted {
			return "D "
		}
		return "M "
	}

	// When not selected, use colored styles
	if isStaged {
		return StatusStagedStyle.Render("S ")
	}
	if node.File.IsNew {
		return StatusNewStyle.Render("??")
	} else if node.File.IsDeleted {
		return StatusDeletedStyle.Render("D ")
	}
	return StatusModifiedStyle.Render("M ")
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

// toggleStaging toggles the staging status of the selected file
func (m *Model) toggleStaging() {
	if len(m.visibleNodes) == 0 || m.selectedIdx >= len(m.visibleNodes) {
		return
	}

	node := m.visibleNodes[m.selectedIdx]
	if node.File == nil || m.gitRunner == nil {
		return
	}

	filepath := node.File.Name
	ctx := context.Background()

	if m.stagedFiles[filepath] {
		// Unstage the file
		err := m.gitRunner.UnstageFile(ctx, filepath)
		if err == nil {
			delete(m.stagedFiles, filepath)
		}
	} else {
		// Stage the file
		err := m.gitRunner.StageFile(ctx, filepath)
		if err == nil {
			m.stagedFiles[filepath] = true
		}
	}
}

// hasStagedFiles returns true if there are any staged files
func (m *Model) hasStagedFiles() bool {
	return len(m.stagedFiles) > 0
}

// openCommitModal opens the commit message modal
func (m *Model) openCommitModal() {
	m.commitModalActive = true
	m.commitError = ""
	m.commitInput.SetValue("")
	m.commitInput.Focus()
}

// closeCommitModal closes the commit message modal
func (m *Model) closeCommitModal() {
	m.commitModalActive = false
	m.commitError = ""
	m.commitInput.Blur()
}

// executeCommit executes the git commit with the entered message
func (m *Model) executeCommit() {
	if m.gitRunner == nil {
		m.commitError = "Git runner not available"
		return
	}

	message := m.commitInput.Value()
	if message == "" {
		m.commitError = "Commit message cannot be empty"
		return
	}

	ctx := context.Background()
	err := m.gitRunner.Commit(ctx, message)
	if err != nil {
		m.commitError = err.Error()
		return
	}

	// Close modal and refresh diff
	m.closeCommitModal()
	m.refreshDiff()
}

// refreshDiff re-runs git diff and rebuilds the file tree
func (m *Model) refreshDiff() {
	if m.gitRunner == nil {
		return
	}

	ctx := context.Background()

	// Re-run git diff
	diffOutput, err := m.gitRunner.RunDiff(ctx, m.diffArgs...)
	if err != nil {
		return
	}

	// Re-parse the diff
	result, err := parser.ParseString(diffOutput)
	if err != nil {
		return
	}

	// Update the model with new files
	m.files = result.Files
	m.treeRoots = BuildTree(m.files, m.rootName)
	m.visibleNodes = FlattenVisible(m.treeRoots)

	// Reset selection to first file
	m.selectedIdx = 0
	for i, node := range m.visibleNodes {
		if node.IsFile() {
			m.selectedIdx = i
			break
		}
	}

	// Reload staged files
	m.stagedFiles = make(map[string]bool)
	staged, err := m.gitRunner.GetStagedFiles(ctx)
	if err == nil {
		for _, f := range staged {
			m.stagedFiles[f] = true
		}
	}

	// Update diff content
	m.updateDiffContent()
}

// renderCommitModal renders the commit message modal overlay
func (m Model) renderCommitModal(background string) string {
	// Modal title
	title := ModalTitleStyle.Render("Commit Message")

	// Text input
	input := m.commitInput.View()

	// Error message if any
	var errorMsg string
	if m.commitError != "" {
		errorMsg = ModalErrorStyle.Render(m.commitError)
	}

	// Help text
	help := ModalHelpStyle.Render("Enter: commit | Esc: cancel")

	// Build modal content
	var modalContent string
	if errorMsg != "" {
		modalContent = lipgloss.JoinVertical(lipgloss.Left, title, input, errorMsg, help)
	} else {
		modalContent = lipgloss.JoinVertical(lipgloss.Left, title, input, help)
	}

	// Style and size the modal
	modal := ModalStyle.Width(60).Render(modalContent)

	// Center modal on screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
}
