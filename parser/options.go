package parser

// Options configures parser behavior
type Options struct {
	// GitPath is the path to the git binary (default: "git")
	GitPath string

	// WorkDir is the working directory for git commands (default: current directory)
	WorkDir string
}

// Option is a functional option for configuring the parser
type Option func(*Options)

// DefaultOptions returns the default parser options
func DefaultOptions() Options {
	return Options{
		GitPath: "git",
		WorkDir: "",
	}
}

// WithGitPath sets the git binary path
func WithGitPath(path string) Option {
	return func(o *Options) {
		o.GitPath = path
	}
}

// WithWorkDir sets the working directory for git commands
func WithWorkDir(dir string) Option {
	return func(o *Options) {
		o.WorkDir = dir
	}
}
