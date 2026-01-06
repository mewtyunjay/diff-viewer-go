package diff

type LineType int

const (
	Context LineType = iota
	Add
	Delete
	Placeholder
)

// Segment represents a portion of a line with optional highlighting
type Segment struct {
	Text        string
	Highlighted bool // true = changed portion, needs intense highlight
}

type Line struct {
	Type     LineType  // what type of line it is
	Content  string    // actual content of that line
	Segments []Segment // nil for non-modified lines; populated for word-level diff
}

type FileDiff struct {
	Name       string
	OldPath    string
	NewPath    string
	IsNew      bool
	IsDeleted  bool
	IsBinary   bool
	AddCount   int
	DelCount   int
	LeftLines  []Line
	RightLines []Line
}

type Result struct {
	Files []FileDiff
}
