package parser

import (
	"os"
	"path/filepath"
	"testing"

	"diff-tui/diff"
)

func TestParseString_BasicDiff(t *testing.T) {
	input := `diff --git a/main.go b/main.go
index 1234567..abcdef0 100644
--- a/main.go
+++ b/main.go
@@ -1,5 +1,6 @@
 package main

 func main() {
-    println("hello")
+    println("hello world")
+    println("new line")
 }
`
	p := New()
	result, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(result.Files))
	}

	file := result.Files[0]
	if file.Name != "main.go" {
		t.Errorf("expected file name 'main.go', got '%s'", file.Name)
	}
	if file.AddCount != 2 {
		t.Errorf("expected 2 additions, got %d", file.AddCount)
	}
	if file.DelCount != 1 {
		t.Errorf("expected 1 deletion, got %d", file.DelCount)
	}
	if file.IsNew {
		t.Error("file should not be marked as new")
	}
	if file.IsDeleted {
		t.Error("file should not be marked as deleted")
	}
}

func TestParseString_NewFile(t *testing.T) {
	input := `diff --git a/utils.go b/utils.go
new file mode 100644
--- /dev/null
+++ b/utils.go
@@ -0,0 +1,5 @@
+package main
+
+func add(a, b int) int {
+    return a + b
+}
`
	p := New()
	result, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(result.Files))
	}

	file := result.Files[0]
	if !file.IsNew {
		t.Error("file should be marked as new")
	}
	if file.AddCount != 5 {
		t.Errorf("expected 5 additions, got %d", file.AddCount)
	}
	if file.DelCount != 0 {
		t.Errorf("expected 0 deletions, got %d", file.DelCount)
	}
}

func TestParseString_DeletedFile(t *testing.T) {
	input := `diff --git a/old.go b/old.go
deleted file mode 100644
--- a/old.go
+++ /dev/null
@@ -1,3 +0,0 @@
-package main
-
-func deprecated() {}
`
	p := New()
	result, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(result.Files))
	}

	file := result.Files[0]
	if !file.IsDeleted {
		t.Error("file should be marked as deleted")
	}
	if file.AddCount != 0 {
		t.Errorf("expected 0 additions, got %d", file.AddCount)
	}
	if file.DelCount != 3 {
		t.Errorf("expected 3 deletions, got %d", file.DelCount)
	}
}

func TestParseString_MultipleFiles(t *testing.T) {
	// Read from testdata
	testdataPath := filepath.Join("..", "testdata", "test.diff")
	data, err := os.ReadFile(testdataPath)
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}

	p := New()
	result, err := p.ParseString(string(data))
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(result.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(result.Files))
	}

	// First file: main.go (modification)
	if result.Files[0].Name != "main.go" {
		t.Errorf("expected first file 'main.go', got '%s'", result.Files[0].Name)
	}
	if result.Files[0].IsNew {
		t.Error("main.go should not be new")
	}

	// Second file: utils.go (new file)
	if result.Files[1].Name != "utils.go" {
		t.Errorf("expected second file 'utils.go', got '%s'", result.Files[1].Name)
	}
	if !result.Files[1].IsNew {
		t.Error("utils.go should be marked as new")
	}
}

func TestParseString_EmptyInput(t *testing.T) {
	p := New()
	_, err := p.ParseString("")
	if err != ErrEmptyDiff {
		t.Errorf("expected ErrEmptyDiff, got %v", err)
	}
}

func TestParseString_WhitespaceOnly(t *testing.T) {
	p := New()
	_, err := p.ParseString("   \n\n   \n")
	if err != ErrEmptyDiff {
		t.Errorf("expected ErrEmptyDiff, got %v", err)
	}
}

func TestParseString_PlaceholderLines(t *testing.T) {
	input := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,3 +1,5 @@
 package main
-func old() {}
+func new1() {}
+func new2() {}
+func new3() {}
`
	p := New()
	result, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	file := result.Files[0]

	// Count placeholder lines
	leftPlaceholders := 0
	rightPlaceholders := 0
	for _, l := range file.LeftLines {
		if l.Type == diff.Placeholder {
			leftPlaceholders++
		}
	}
	for _, l := range file.RightLines {
		if l.Type == diff.Placeholder {
			rightPlaceholders++
		}
	}

	// 1 delete, 3 adds -> 2 extra adds need placeholders on left
	if leftPlaceholders != 2 {
		t.Errorf("expected 2 left placeholders, got %d", leftPlaceholders)
	}
	if rightPlaceholders != 0 {
		t.Errorf("expected 0 right placeholders, got %d", rightPlaceholders)
	}
}
