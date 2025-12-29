package tui

import "diff-tui/diff"

// MockFile represents a file with diff data for display
type MockFile struct {
	Name      string
	AddCount  int
	DelCount  int
	IsNew     bool
	IsDeleted bool
	LeftLines  []diff.Line  // Original side (context + deletions)
	RightLines []diff.Line  // Modified side (context + additions)
}

// GetMockFiles returns sample diff data for UI demonstration
func GetMockFiles() []MockFile {
	return []MockFile{
		{
			Name:     "main.go",
			AddCount: 2,
			DelCount: 1,
			LeftLines: []diff.Line{
				{Type: diff.Context, Content: "package main"},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: "import \"fmt\""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: "func main() {"},
				{Type: diff.Delete, Content: "    fmt.Println(\"hello\")"},
				{Type: diff.Context, Content: ""},  // placeholder for add
				{Type: diff.Context, Content: ""},  // placeholder for add
				{Type: diff.Context, Content: "}"},
			},
			RightLines: []diff.Line{
				{Type: diff.Context, Content: "package main"},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: "import \"fmt\""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: "func main() {"},
				{Type: diff.Context, Content: ""},  // placeholder for delete
				{Type: diff.Add, Content: "    fmt.Println(\"hello world\")"},
				{Type: diff.Add, Content: "    fmt.Println(\"new line\")"},
				{Type: diff.Context, Content: "}"},
			},
		},
		{
			Name:     "utils.go",
			AddCount: 5,
			DelCount: 0,
			IsNew:    true,
			LeftLines: []diff.Line{
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
			},
			RightLines: []diff.Line{
				{Type: diff.Add, Content: "package main"},
				{Type: diff.Add, Content: ""},
				{Type: diff.Add, Content: "func add(a, b int) int {"},
				{Type: diff.Add, Content: "    return a + b"},
				{Type: diff.Add, Content: "}"},
			},
		},
		{
			Name:     "config.go",
			AddCount: 3,
			DelCount: 2,
			LeftLines: []diff.Line{
				{Type: diff.Context, Content: "package main"},
				{Type: diff.Context, Content: ""},
				{Type: diff.Delete, Content: "const Version = \"1.0.0\""},
				{Type: diff.Delete, Content: "const Debug = false"},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: "type Config struct {"},
				{Type: diff.Context, Content: "    Name string"},
				{Type: diff.Context, Content: "}"},
			},
			RightLines: []diff.Line{
				{Type: diff.Context, Content: "package main"},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Add, Content: "const Version = \"2.0.0\""},
				{Type: diff.Add, Content: "const Debug = true"},
				{Type: diff.Add, Content: "const MaxRetries = 3"},
				{Type: diff.Context, Content: "type Config struct {"},
				{Type: diff.Context, Content: "    Name string"},
				{Type: diff.Context, Content: "}"},
			},
		},
		{
			Name:      "old_utils.go",
			AddCount:  0,
			DelCount:  4,
			IsDeleted: true,
			LeftLines: []diff.Line{
				{Type: diff.Delete, Content: "package main"},
				{Type: diff.Delete, Content: ""},
				{Type: diff.Delete, Content: "func deprecated() {"},
				{Type: diff.Delete, Content: "}"},
			},
			RightLines: []diff.Line{
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
				{Type: diff.Context, Content: ""},
			},
		},
	}
}
