package diff

type LineType int

const (
	Context LineType = iota
	Add
	Delete
)

// a single line in a hunk
type Line struct {
	Type    LineType // what type of line it is
	Content string   // actual content of that line
}
