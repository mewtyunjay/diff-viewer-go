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

// StageFile stages a single file (git add)
func (g *GitRunner) StageFile(ctx context.Context, filepath string) error {
	cmd := exec.CommandContext(ctx, g.gitPath, "add", filepath)
	if g.workDir != "" {
		cmd.Dir = g.workDir
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return &GitError{
			Args:   []string{"add", filepath},
			Stderr: strings.TrimSpace(stderr.String()),
			Err:    err,
		}
	}

	return nil
}

// UnstageFile unstages a single file (git reset HEAD)
func (g *GitRunner) UnstageFile(ctx context.Context, filepath string) error {
	cmd := exec.CommandContext(ctx, g.gitPath, "reset", "HEAD", filepath)
	if g.workDir != "" {
		cmd.Dir = g.workDir
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return &GitError{
			Args:   []string{"reset", "HEAD", filepath},
			Stderr: strings.TrimSpace(stderr.String()),
			Err:    err,
		}
	}

	return nil
}

// Commit creates a commit with the given message
func (g *GitRunner) Commit(ctx context.Context, message string) error {
	cmd := exec.CommandContext(ctx, g.gitPath, "commit", "-m", message)
	if g.workDir != "" {
		cmd.Dir = g.workDir
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return &GitError{
			Args:   []string{"commit", "-m", message},
			Stderr: strings.TrimSpace(stderr.String()),
			Err:    err,
		}
	}

	return nil
}

// GetStagedFiles returns a list of currently staged file paths
func (g *GitRunner) GetStagedFiles(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, g.gitPath, "diff", "--cached", "--name-only")
	if g.workDir != "" {
		cmd.Dir = g.workDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, &GitError{
			Args:   []string{"diff", "--cached", "--name-only"},
			Stderr: strings.TrimSpace(stderr.String()),
			Err:    err,
		}
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return nil, nil
	}

	return strings.Split(output, "\n"), nil
}
