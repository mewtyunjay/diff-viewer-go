package parser

import "diff-tui/diff"

// alignHunks transforms parsed hunks into aligned left/right line slices
// Returns: leftLines, rightLines, addCount, delCount
func alignHunks(hunks []hunk) ([]diff.Line, []diff.Line, int, int) {
	var left, right []diff.Line
	var addCount, delCount int

	for _, h := range hunks {
		l, r, adds, dels := alignHunk(h)
		left = append(left, l...)
		right = append(right, r...)
		addCount += adds
		delCount += dels
	}

	return left, right, addCount, delCount
}

// alignHunk aligns a single hunk for side-by-side display
func alignHunk(h hunk) ([]diff.Line, []diff.Line, int, int) {
	var left, right []diff.Line
	var addCount, delCount int

	i := 0
	for i < len(h.lines) {
		line := h.lines[i]

		switch line.Type {
		case diff.Context:
			// Context lines appear on both sides
			left = append(left, line)
			right = append(right, line)
			i++

		case diff.Delete:
			// Collect consecutive deletes
			var deletes []diff.Line
			for i < len(h.lines) && h.lines[i].Type == diff.Delete {
				deletes = append(deletes, h.lines[i])
				delCount++
				i++
			}

			// Collect consecutive adds that follow
			var adds []diff.Line
			for i < len(h.lines) && h.lines[i].Type == diff.Add {
				adds = append(adds, h.lines[i])
				addCount++
				i++
			}

			// Align deletes and adds side by side
			alignBlock(deletes, adds, &left, &right)

		case diff.Add:
			// Standalone adds (not following deletes)
			var adds []diff.Line
			for i < len(h.lines) && h.lines[i].Type == diff.Add {
				adds = append(adds, h.lines[i])
				addCount++
				i++
			}

			// Add placeholders on left, adds on right
			for _, add := range adds {
				left = append(left, diff.Line{Type: diff.Placeholder, Content: ""})
				right = append(right, add)
			}
		}
	}

	return left, right, addCount, delCount
}

// alignBlock aligns a block of deletes and adds side by side
// For modifications (delete followed by add), they appear on the same row
// and word-level diff highlighting is computed for paired lines.
func alignBlock(deletes, adds []diff.Line, left, right *[]diff.Line) {
	maxLen := len(deletes)
	if len(adds) > maxLen {
		maxLen = len(adds)
	}

	for i := 0; i < maxLen; i++ {
		var leftLine, rightLine diff.Line

		if i < len(deletes) {
			leftLine = deletes[i]
		} else {
			// Placeholder for add-only line
			leftLine = diff.Line{Type: diff.Placeholder, Content: ""}
		}

		if i < len(adds) {
			rightLine = adds[i]
		} else {
			// Placeholder for delete-only line
			rightLine = diff.Line{Type: diff.Placeholder, Content: ""}
		}

		// Compute word-level diff for modification pairs (delete + add on same row)
		if leftLine.Type == diff.Delete && rightLine.Type == diff.Add {
			leftSegs, rightSegs := diff.ComputeWordDiff(leftLine.Content, rightLine.Content)
			leftLine.Segments = leftSegs
			rightLine.Segments = rightSegs
		}

		*left = append(*left, leftLine)
		*right = append(*right, rightLine)
	}
}
