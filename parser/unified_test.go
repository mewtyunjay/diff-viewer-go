package parser

import (
	"testing"

	"diff-tui/diff"
)

func TestParseUnified_HunkHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		oldStart int
		oldCount int
		newStart int
		newCount int
	}{
		{
			name:     "standard hunk",
			input:    "@@ -1,5 +1,6 @@",
			oldStart: 1, oldCount: 5, newStart: 1, newCount: 6,
		},
		{
			name:     "single line old",
			input:    "@@ -1 +1,3 @@",
			oldStart: 1, oldCount: 1, newStart: 1, newCount: 3,
		},
		{
			name:     "single line new",
			input:    "@@ -1,3 +1 @@",
			oldStart: 1, oldCount: 3, newStart: 1, newCount: 1,
		},
		{
			name:     "both single line",
			input:    "@@ -5 +10 @@",
			oldStart: 5, oldCount: 1, newStart: 10, newCount: 1,
		},
		{
			name:     "new file",
			input:    "@@ -0,0 +1,5 @@",
			oldStart: 0, oldCount: 0, newStart: 1, newCount: 5,
		},
		{
			name:     "deleted file",
			input:    "@@ -1,5 +0,0 @@",
			oldStart: 1, oldCount: 5, newStart: 0, newCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := hunkHeaderRE.FindStringSubmatch(tt.input)
			if matches == nil {
				t.Fatal("hunk header regex did not match")
			}

			// Parse values manually like parseHunk does
			hp := &HunkParser{lines: []string{tt.input}, pos: 0}
			h, _, err := hp.parseHunk()
			if err != nil {
				t.Fatalf("parseHunk failed: %v", err)
			}

			if h.oldStart != tt.oldStart {
				t.Errorf("oldStart: expected %d, got %d", tt.oldStart, h.oldStart)
			}
			if h.oldCount != tt.oldCount {
				t.Errorf("oldCount: expected %d, got %d", tt.oldCount, h.oldCount)
			}
			if h.newStart != tt.newStart {
				t.Errorf("newStart: expected %d, got %d", tt.newStart, h.newStart)
			}
			if h.newCount != tt.newCount {
				t.Errorf("newCount: expected %d, got %d", tt.newCount, h.newCount)
			}
		})
	}
}

// HunkParser is a helper for testing parseHunk directly
type HunkParser struct {
	lines []string
	pos   int
}

func (hp *HunkParser) parseHunk() (hunk, int, error) {
	return parseHunk(hp.lines, hp.pos)
}

func TestParseUnified_LineTypes(t *testing.T) {
	input := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,4 +1,4 @@
 context line
-deleted line
+added line
 another context
`
	files, err := parseUnified(input)
	if err != nil {
		t.Fatalf("parseUnified failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	// Check that we have the right mix of line types
	hasContext := false
	hasAdd := false
	hasDelete := false
	hasPlaceholder := false

	for _, l := range files[0].LeftLines {
		switch l.Type {
		case diff.Context:
			hasContext = true
		case diff.Delete:
			hasDelete = true
		case diff.Placeholder:
			hasPlaceholder = true
		}
	}
	for _, l := range files[0].RightLines {
		if l.Type == diff.Add {
			hasAdd = true
		}
	}

	if !hasContext {
		t.Error("expected context lines")
	}
	if !hasAdd {
		t.Error("expected add lines")
	}
	if !hasDelete {
		t.Error("expected delete lines")
	}
	// In this case, 1 delete + 1 add = no placeholders needed
	if hasPlaceholder {
		t.Error("did not expect placeholders for equal add/delete count")
	}
}

func TestParseUnified_BinaryFile(t *testing.T) {
	input := `diff --git a/image.png b/image.png
new file mode 100644
Binary files /dev/null and b/image.png differ
`
	files, err := parseUnified(input)
	if err != nil {
		t.Fatalf("parseUnified failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	if !files[0].IsBinary {
		t.Error("file should be marked as binary")
	}
	if !files[0].IsNew {
		t.Error("file should be marked as new")
	}
}

func TestParseUnified_FileHeaders(t *testing.T) {
	input := `diff --git a/old/path.go b/new/path.go
--- a/old/path.go
+++ b/new/path.go
@@ -1 +1 @@
-old
+new
`
	files, err := parseUnified(input)
	if err != nil {
		t.Fatalf("parseUnified failed: %v", err)
	}

	file := files[0]
	if file.OldPath != "old/path.go" {
		t.Errorf("expected OldPath 'old/path.go', got '%s'", file.OldPath)
	}
	if file.NewPath != "new/path.go" {
		t.Errorf("expected NewPath 'new/path.go', got '%s'", file.NewPath)
	}
	// Name should be the new path
	if file.Name != "new/path.go" {
		t.Errorf("expected Name 'new/path.go', got '%s'", file.Name)
	}
}

func TestParseUnified_NoNewlineAtEOF(t *testing.T) {
	input := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1 +1 @@
-old
\ No newline at end of file
+new
\ No newline at end of file
`
	files, err := parseUnified(input)
	if err != nil {
		t.Fatalf("parseUnified failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	// Should have parsed correctly despite the backslash lines
	if files[0].DelCount != 1 {
		t.Errorf("expected 1 deletion, got %d", files[0].DelCount)
	}
	if files[0].AddCount != 1 {
		t.Errorf("expected 1 addition, got %d", files[0].AddCount)
	}
}

func TestParseUnified_EmptyHunk(t *testing.T) {
	// This is an edge case - a hunk with no actual changes
	input := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,2 +1,2 @@
 line1
 line2
`
	files, err := parseUnified(input)
	if err != nil {
		t.Fatalf("parseUnified failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	// No actual changes, just context
	if files[0].AddCount != 0 {
		t.Errorf("expected 0 additions, got %d", files[0].AddCount)
	}
	if files[0].DelCount != 0 {
		t.Errorf("expected 0 deletions, got %d", files[0].DelCount)
	}
}
