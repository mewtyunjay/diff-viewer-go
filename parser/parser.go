package parser

import (
	"context"
	"io"
	"strings"

	"diff-tui/diff"
)

// Parser is the main diff parser
type Parser struct {
	opts Options
	git  *GitRunner
}

// New creates a new Parser with the given options
func New(opts ...Option) *Parser {
	o := DefaultOptions()
	for _, opt := range opts {
		opt(&o)
	}

	return &Parser{
		opts: o,
		git:  NewGitRunner(o.GitPath, o.WorkDir),
	}
}

// ParseString parses a unified diff from a string (method on Parser)
func (p *Parser) ParseString(input string) (*diff.Result, error) {
	return ParseString(input)
}

// ParseString parses a unified diff from a string (standalone function)
func ParseString(input string) (*diff.Result, error) {
	files, err := parseUnified(input)
	if err != nil {
		return nil, err
	}

	return &diff.Result{Files: files}, nil
}

// ParseReader parses a unified diff from an io.Reader
func (p *Parser) ParseReader(r io.Reader) (*diff.Result, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return p.ParseString(string(data))
}

// ParseGitDiff executes git diff with the given args and parses the output
func (p *Parser) ParseGitDiff(ctx context.Context, args ...string) (*diff.Result, error) {
	// Check if we're in a git repo
	if !p.git.IsGitRepository(ctx) {
		return nil, ErrNotGitRepo
	}

	// Run git diff
	output, err := p.git.RunDiff(ctx, args...)
	if err != nil {
		return nil, err
	}

	// Handle empty diff
	if strings.TrimSpace(output) == "" {
		return nil, ErrEmptyDiff
	}

	return p.ParseString(output)
}

// IsGitRepository checks if the working directory is a git repository
func (p *Parser) IsGitRepository(ctx context.Context) bool {
	return p.git.IsGitRepository(ctx)
}

// GitRunner returns the underlying git runner for direct git operations
func (p *Parser) GitRunner() *GitRunner {
	return p.git
}
