package tui

import (
	"sort"
	"strings"

	"diff-tui/diff"
)

// NodeType represents the type of tree node
type NodeType int

const (
	NodeDirectory NodeType = iota
	NodeFile
)

// TreeNode represents a node in the file tree
type TreeNode struct {
	Name     string          // Just the filename/dirname, not full path
	Path     string          // Full path for lookup
	Type     NodeType
	Children []*TreeNode
	Expanded bool            // For directories (default: true)
	File     *diff.FileDiff  // nil for directories
	Depth    int             // Indentation level
	Parent   *TreeNode       // For navigation (go to parent)
}

// BuildTree creates a tree structure from a flat list of files
func BuildTree(files []diff.FileDiff) []*TreeNode {
	if len(files) == 0 {
		return nil
	}

	// Map to track created directories
	dirNodes := make(map[string]*TreeNode)
	var roots []*TreeNode

	for i := range files {
		file := &files[i]
		path := file.Name

		// Split path into parts
		parts := strings.Split(path, "/")

		var parent *TreeNode
		currentPath := ""

		// Create directory nodes as needed
		for j, part := range parts {
			if j < len(parts)-1 {
				// This is a directory
				if currentPath == "" {
					currentPath = part
				} else {
					currentPath = currentPath + "/" + part
				}

				if _, exists := dirNodes[currentPath]; !exists {
					node := &TreeNode{
						Name:     part,
						Path:     currentPath,
						Type:     NodeDirectory,
						Expanded: true, // Start expanded by default
						Parent:   parent,
						Depth:    j,
					}
					dirNodes[currentPath] = node

					if parent == nil {
						roots = append(roots, node)
					} else {
						parent.Children = append(parent.Children, node)
					}
				}
				parent = dirNodes[currentPath]
			} else {
				// This is the file
				node := &TreeNode{
					Name:   part,
					Path:   path,
					Type:   NodeFile,
					File:   file,
					Parent: parent,
					Depth:  j,
				}

				if parent == nil {
					roots = append(roots, node)
				} else {
					parent.Children = append(parent.Children, node)
				}
			}
		}
	}

	// Sort all children: directories first, then alphabetically
	sortChildren(roots)
	for _, dir := range dirNodes {
		sortChildren(dir.Children)
	}

	return roots
}

// sortChildren sorts nodes: directories first, then files, alphabetically within each group
func sortChildren(nodes []*TreeNode) {
	sort.Slice(nodes, func(i, j int) bool {
		// Directories come before files
		if nodes[i].Type != nodes[j].Type {
			return nodes[i].Type == NodeDirectory
		}
		// Alphabetically within same type
		return strings.ToLower(nodes[i].Name) < strings.ToLower(nodes[j].Name)
	})
}

// FlattenVisible returns a flat list of visible nodes for navigation
func FlattenVisible(roots []*TreeNode) []*TreeNode {
	var result []*TreeNode
	for _, root := range roots {
		flattenNode(root, &result)
	}
	return result
}

func flattenNode(node *TreeNode, result *[]*TreeNode) {
	*result = append(*result, node)

	if node.Type == NodeDirectory && node.Expanded {
		for _, child := range node.Children {
			flattenNode(child, result)
		}
	}
}

// ToggleExpanded toggles the expanded state of a directory node
func (n *TreeNode) ToggleExpanded() {
	if n.Type == NodeDirectory {
		n.Expanded = !n.Expanded
	}
}

// IsFile returns true if the node is a file
func (n *TreeNode) IsFile() bool {
	return n.Type == NodeFile
}

// IsDirectory returns true if the node is a directory
func (n *TreeNode) IsDirectory() bool {
	return n.Type == NodeDirectory
}


// FindFirstFile finds the first file node in the tree (for initial selection)
func FindFirstFile(roots []*TreeNode) *TreeNode {
	visible := FlattenVisible(roots)
	for _, node := range visible {
		if node.IsFile() {
			return node
		}
	}
	// If no files, return first node
	if len(visible) > 0 {
		return visible[0]
	}
	return nil
}
