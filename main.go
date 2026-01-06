package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"diff-tui/parser"
	"diff-tui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx := context.Background()
	args := os.Args[1:]

	// Parse git diff with provided arguments
	p := parser.New()
	result, err := p.ParseGitDiff(ctx, args...)
	if err != nil {
		handleError(err)
		return
	}

	// Get the git root folder name for display in the tree
	rootName := ""
	if gitRoot, err := p.GitRunner().FindGitRoot(ctx); err == nil {
		rootName = filepath.Base(gitRoot)
	}

	// Pass GitRunner, args, and rootName to enable staging/commit features
	model := tui.NewModel(result.Files, p.GitRunner(), args, rootName)
	runTUI(model)
}

func runTUI(model tui.Model) {
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleError(err error) {
	switch {
	case errors.Is(err, parser.ErrNotGitRepo):
		fmt.Fprintln(os.Stderr, "Error: Not a git repository")
		fmt.Fprintln(os.Stderr, "Run this command from within a git repository")
		os.Exit(1)

	case errors.Is(err, parser.ErrEmptyDiff):
		fmt.Fprintln(os.Stderr, "No changes to display")
		fmt.Fprintln(os.Stderr, "Try: diff-tui HEAD~1  or  diff-tui --staged")
		os.Exit(0)

	default:
		var gitErr *parser.GitError
		if errors.As(err, &gitErr) {
			fmt.Fprintf(os.Stderr, "Git error: %s\n", gitErr.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
