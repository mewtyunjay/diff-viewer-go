package parser

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitRunner handles git command execution
type GitRunner struct {
	gitPath string
	workDir string
}

// NewGitRunner creates a new git runner
func NewGitRunner(gitPath, workDir string) *GitRunner {
	if gitPath == "" {
		gitPath = "git"
	}
	return &GitRunner{
		gitPath: gitPath,
		workDir: workDir,
	}
}

// IsGitRepository checks if workDir is inside a git repository
func (g *GitRunner) IsGitRepository(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, g.gitPath, "rev-parse", "--git-dir")
	if g.workDir != "" {
		cmd.Dir = g.workDir
	}
	err := cmd.Run()
	return err == nil
}

// RunDiff executes git diff with the given arguments
func (g *GitRunner) RunDiff(ctx context.Context, args ...string) (string, error) {
	// Build command args: always include --no-color to avoid ANSI codes
	cmdArgs := []string{"diff", "--no-color"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.CommandContext(ctx, g.gitPath, cmdArgs...)
	if g.workDir != "" {
		cmd.Dir = g.workDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Check if it's a "not a git repo" error
		if strings.Contains(stderr.String(), "not a git repository") {
			return "", ErrNotGitRepo
		}
		return "", &GitError{
			Args:   cmdArgs,
			Stderr: strings.TrimSpace(stderr.String()),
			Err:    err,
		}
	}

	return stdout.String(), nil
}

// FindGitRoot finds the root directory of the git repository
func (g *GitRunner) FindGitRoot(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, g.gitPath, "rev-parse", "--show-toplevel")
	if g.workDir != "" {
		cmd.Dir = g.workDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if strings.Contains(stderr.String(), "not a git repository") {
			return "", ErrNotGitRepo
		}
		return "", &GitError{
			Args:   []string{"rev-parse", "--show-toplevel"},
			Stderr: strings.TrimSpace(stderr.String()),
			Err:    err,
		}
	}

	return filepath.Clean(strings.TrimSpace(stdout.String())), nil
}
