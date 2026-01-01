package tui

import (
	"testing"

	"diff-tui/diff"
)

func TestBuildTree_SingleFile(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "main.go", AddCount: 5, DelCount: 3},
	}

	roots := BuildTree(files)

	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}

	if roots[0].Name != "main.go" {
		t.Errorf("expected name 'main.go', got '%s'", roots[0].Name)
	}

	if roots[0].Type != NodeFile {
		t.Error("expected NodeFile type")
	}
}

func TestBuildTree_NestedFiles(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "src/utils/helper.go", AddCount: 10, DelCount: 0},
		{Name: "src/main.go", AddCount: 5, DelCount: 3},
	}

	roots := BuildTree(files)

	if len(roots) != 1 {
		t.Fatalf("expected 1 root (src/), got %d", len(roots))
	}

	src := roots[0]
	if src.Name != "src" {
		t.Errorf("expected root name 'src', got '%s'", src.Name)
	}

	if src.Type != NodeDirectory {
		t.Error("expected NodeDirectory type for src")
	}

	if !src.Expanded {
		t.Error("expected directory to be expanded by default")
	}

	// Should have 2 children: main.go and utils/
	if len(src.Children) != 2 {
		t.Fatalf("expected 2 children in src, got %d", len(src.Children))
	}
}

func TestBuildTree_MultipleRoots(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "README.md", AddCount: 1, DelCount: 0},
		{Name: "src/main.go", AddCount: 5, DelCount: 3},
	}

	roots := BuildTree(files)

	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}

	// Directories should come first (sorted)
	if roots[0].Name != "src" {
		t.Errorf("expected first root 'src', got '%s'", roots[0].Name)
	}

	if roots[1].Name != "README.md" {
		t.Errorf("expected second root 'README.md', got '%s'", roots[1].Name)
	}
}

func TestFlattenVisible_Expanded(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "src/main.go"},
		{Name: "src/utils.go"},
	}

	roots := BuildTree(files)
	visible := FlattenVisible(roots)

	// Should show: src/, main.go, utils.go
	if len(visible) != 3 {
		t.Fatalf("expected 3 visible nodes, got %d", len(visible))
	}

	if visible[0].Name != "src" {
		t.Errorf("expected first visible 'src', got '%s'", visible[0].Name)
	}
}

func TestFlattenVisible_Collapsed(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "src/main.go"},
		{Name: "src/utils.go"},
	}

	roots := BuildTree(files)
	roots[0].Expanded = false // Collapse src/

	visible := FlattenVisible(roots)

	// Should only show: src/
	if len(visible) != 1 {
		t.Fatalf("expected 1 visible node, got %d", len(visible))
	}

	if visible[0].Name != "src" {
		t.Errorf("expected visible 'src', got '%s'", visible[0].Name)
	}
}

func TestTreeNode_Parent(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "src/main.go"},
	}

	roots := BuildTree(files)

	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}

	src := roots[0]
	if len(src.Children) != 1 {
		t.Fatalf("expected 1 child in src, got %d", len(src.Children))
	}

	mainGo := src.Children[0]
	if mainGo.Parent != src {
		t.Error("expected main.go's parent to be src")
	}

	if src.Parent != nil {
		t.Error("expected src's parent to be nil")
	}
}

func TestTreeNode_ToggleExpanded(t *testing.T) {
	node := &TreeNode{
		Type:     NodeDirectory,
		Expanded: true,
	}

	node.ToggleExpanded()
	if node.Expanded {
		t.Error("expected Expanded to be false after toggle")
	}

	node.ToggleExpanded()
	if !node.Expanded {
		t.Error("expected Expanded to be true after second toggle")
	}
}

func TestFindFirstFile(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "src/main.go"},
	}

	roots := BuildTree(files)
	first := FindFirstFile(roots)

	if first == nil {
		t.Fatal("expected to find first file")
	}

	if first.Name != "main.go" {
		t.Errorf("expected 'main.go', got '%s'", first.Name)
	}
}
