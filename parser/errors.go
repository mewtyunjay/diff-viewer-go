package parser

import (
	"errors"
	"fmt"
)

var (
	// ErrNotGitRepo indicates the directory is not a git repository
	ErrNotGitRepo = errors.New("not a git repository")

	// ErrEmptyDiff indicates no changes to display
	ErrEmptyDiff = errors.New("no diff output")

	// ErrInvalidDiff indicates malformed diff input
	ErrInvalidDiff = errors.New("invalid diff format")
)

// GitError wraps errors from git command execution
type GitError struct {
	Args   []string
	Stderr string
	Err    error
}

func (e *GitError) Error() string {
	if e.Stderr != "" {
		return fmt.Sprintf("git command failed: %s", e.Stderr)
	}
	return fmt.Sprintf("git command failed: %v", e.Err)
}

func (e *GitError) Unwrap() error {
	return e.Err
}

// ParseError provides detailed error information for parsing failures
type ParseError struct {
	Line    int
	Message string
	Cause   error
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("parse error at line %d: %s", e.Line, e.Message)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Cause
}
