package tui

import (
	"testing"

	"diff-tui/diff"
)

func TestBuildTree_SingleFile(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "main.go", AddCount: 5, DelCount: 3},
	}

	roots := BuildTree(files, "")

	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}

	// Root should be "."
	if roots[0].Name != "." {
		t.Errorf("expected root name '.', got '%s'", roots[0].Name)
	}

	if roots[0].Type != NodeDirectory {
		t.Error("expected NodeDirectory type for root")
	}

	// main.go should be a child of root
	if len(roots[0].Children) != 1 {
		t.Fatalf("expected 1 child in root, got %d", len(roots[0].Children))
	}

	if roots[0].Children[0].Name != "main.go" {
		t.Errorf("expected child name 'main.go', got '%s'", roots[0].Children[0].Name)
	}

	if roots[0].Children[0].Type != NodeFile {
		t.Error("expected NodeFile type for main.go")
	}
}

func TestBuildTree_NestedFiles(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "src/utils/helper.go", AddCount: 10, DelCount: 0},
		{Name: "src/main.go", AddCount: 5, DelCount: 3},
	}

	roots := BuildTree(files, "")

	if len(roots) != 1 {
		t.Fatalf("expected 1 root (.), got %d", len(roots))
	}

	root := roots[0]
	if root.Name != "." {
		t.Errorf("expected root name '.', got '%s'", root.Name)
	}

	// Root should have 1 child: src/
	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child in root, got %d", len(root.Children))
	}

	src := root.Children[0]
	if src.Name != "src" {
		t.Errorf("expected child name 'src', got '%s'", src.Name)
	}

	if src.Type != NodeDirectory {
		t.Error("expected NodeDirectory type for src")
	}

	if !src.Expanded {
		t.Error("expected directory to be expanded by default")
	}

	// src should have 2 children: main.go and utils/
	if len(src.Children) != 2 {
		t.Fatalf("expected 2 children in src, got %d", len(src.Children))
	}
}

func TestBuildTree_MultipleRoots(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "README.md", AddCount: 1, DelCount: 0},
		{Name: "src/main.go", AddCount: 5, DelCount: 3},
	}

	roots := BuildTree(files, "")

	// Now there's only 1 root: "."
	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}

	root := roots[0]
	if root.Name != "." {
		t.Errorf("expected root name '.', got '%s'", root.Name)
	}

	// Root should have 2 children: src/ and README.md
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children in root, got %d", len(root.Children))
	}

	// Directories should come first (sorted)
	if root.Children[0].Name != "src" {
		t.Errorf("expected first child 'src', got '%s'", root.Children[0].Name)
	}

	if root.Children[1].Name != "README.md" {
		t.Errorf("expected second child 'README.md', got '%s'", root.Children[1].Name)
	}
}

func TestFlattenVisible_Expanded(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "src/main.go"},
		{Name: "src/utils.go"},
	}

	roots := BuildTree(files, "")
	visible := FlattenVisible(roots)

	// Should show: ., src/, main.go, utils.go
	if len(visible) != 4 {
		t.Fatalf("expected 4 visible nodes, got %d", len(visible))
	}

	if visible[0].Name != "." {
		t.Errorf("expected first visible '.', got '%s'", visible[0].Name)
	}

	if visible[1].Name != "src" {
		t.Errorf("expected second visible 'src', got '%s'", visible[1].Name)
	}
}

func TestFlattenVisible_Collapsed(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "src/main.go"},
		{Name: "src/utils.go"},
	}

	roots := BuildTree(files, "")
	// Collapse src/ (child of root)
	src := roots[0].Children[0]
	src.Expanded = false

	visible := FlattenVisible(roots)

	// Should show: ., src/
	if len(visible) != 2 {
		t.Fatalf("expected 2 visible nodes, got %d", len(visible))
	}

	if visible[0].Name != "." {
		t.Errorf("expected first visible '.', got '%s'", visible[0].Name)
	}

	if visible[1].Name != "src" {
		t.Errorf("expected second visible 'src', got '%s'", visible[1].Name)
	}
}

func TestTreeNode_Parent(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "src/main.go"},
	}

	roots := BuildTree(files, "")

	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}

	root := roots[0]
	src := root.Children[0]
	if len(src.Children) != 1 {
		t.Fatalf("expected 1 child in src, got %d", len(src.Children))
	}

	mainGo := src.Children[0]
	if mainGo.Parent != src {
		t.Error("expected main.go's parent to be src")
	}

	if src.Parent != root {
		t.Error("expected src's parent to be root")
	}

	if root.Parent != nil {
		t.Error("expected root's parent to be nil")
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

	roots := BuildTree(files, "")
	first := FindFirstFile(roots)

	if first == nil {
		t.Fatal("expected to find first file")
	}

	if first.Name != "main.go" {
		t.Errorf("expected 'main.go', got '%s'", first.Name)
	}
}

func TestBuildTree_CustomRootName(t *testing.T) {
	files := []diff.FileDiff{
		{Name: "main.go"},
	}

	roots := BuildTree(files, "my-project")

	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}

	if roots[0].Name != "my-project" {
		t.Errorf("expected root name 'my-project', got '%s'", roots[0].Name)
	}

	// Path should still be "." for internal consistency
	if roots[0].Path != "." {
		t.Errorf("expected root path '.', got '%s'", roots[0].Path)
	}
}
