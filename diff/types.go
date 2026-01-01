package diff

type LineType int

const (
	Context LineType = iota
	Add
	Delete
	Placeholder
)

type Line struct {
	Type    LineType // what type of line it is
	Content string   // actual content of that line
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
