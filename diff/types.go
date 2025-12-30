package diff

type LineType int

const (
	Context LineType = iota
	Add
	Delete
	Placeholder // Empty filler line for alignment in side-by-side view
)

// Line represents a single line in a hunk
type Line struct {
	Type    LineType // what type of line it is
	Content string   // actual content of that line
}

// FileDiff represents a single file's diff
type FileDiff struct {
	Name       string // display name (usually NewPath or OldPath)
	OldPath    string // path before (--- line)
	NewPath    string // path after (+++ line)
	IsNew      bool   // true if file is newly created
	IsDeleted  bool   // true if file is deleted
	IsBinary   bool   // true if binary file
	AddCount   int    // number of added lines
	DelCount   int    // number of deleted lines
	LeftLines  []Line // aligned lines for left panel (original)
	RightLines []Line // aligned lines for right panel (modified)
}

// Result is the complete parsed diff result
type Result struct {
	Files []FileDiff
}
