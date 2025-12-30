package parser

import (
	"regexp"
	"strconv"
	"strings"

	"diff-tui/diff"
)

var (
	// diffGitRE matches "diff --git a/path b/path"
	diffGitRE = regexp.MustCompile(`^diff --git a/(.+) b/(.+)$`)

	// hunkHeaderRE matches "@@ -start,count +start,count @@" with optional context
	hunkHeaderRE = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)

	// binaryFileRE matches "Binary files ... differ"
	binaryFileRE = regexp.MustCompile(`^Binary files .* differ$`)
)

// hunk represents a parsed hunk before alignment
type hunk struct {
	oldStart int
	oldCount int
	newStart int
	newCount int
	lines    []diff.Line
}

// parseUnified parses unified diff format into FileDiff structs
func parseUnified(input string) ([]diff.FileDiff, error) {
	if input == "" {
		return nil, ErrEmptyDiff
	}

	lines := strings.Split(input, "\n")
	var files []diff.FileDiff
	pos := 0

	for pos < len(lines) {
		// Skip empty lines
		if strings.TrimSpace(lines[pos]) == "" {
			pos++
			continue
		}

		// Look for diff --git header
		if !strings.HasPrefix(lines[pos], "diff --git ") {
			pos++
			continue
		}

		fd, newPos, err := parseFileDiff(lines, pos)
		if err != nil {
			return nil, err
		}
		if fd != nil {
			files = append(files, *fd)
		}
		pos = newPos
	}

	if len(files) == 0 {
		return nil, ErrEmptyDiff
	}

	return files, nil
}

// parseFileDiff parses a single file's diff starting at pos
func parseFileDiff(lines []string, pos int) (*diff.FileDiff, int, error) {
	fd := &diff.FileDiff{}

	// Parse diff --git header
	matches := diffGitRE.FindStringSubmatch(lines[pos])
	if matches == nil {
		return nil, pos + 1, nil
	}
	fd.OldPath = matches[1]
	fd.NewPath = matches[2]
	fd.Name = fd.NewPath
	pos++

	// Parse optional headers (old mode, new mode, new file, deleted file, index, etc.)
	for pos < len(lines) {
		line := lines[pos]

		if strings.HasPrefix(line, "old mode ") {
			pos++
			continue
		}
		if strings.HasPrefix(line, "new mode ") {
			pos++
			continue
		}
		if strings.HasPrefix(line, "new file mode ") {
			fd.IsNew = true
			pos++
			continue
		}
		if strings.HasPrefix(line, "deleted file mode ") {
			fd.IsDeleted = true
			pos++
			continue
		}
		if strings.HasPrefix(line, "index ") {
			pos++
			continue
		}
		if strings.HasPrefix(line, "similarity index ") {
			pos++
			continue
		}
		if strings.HasPrefix(line, "rename from ") {
			pos++
			continue
		}
		if strings.HasPrefix(line, "rename to ") {
			pos++
			continue
		}
		if binaryFileRE.MatchString(line) {
			fd.IsBinary = true
			pos++
			continue
		}

		// Check for --- line
		if strings.HasPrefix(line, "--- ") {
			break
		}

		// Check for next file diff or end
		if strings.HasPrefix(line, "diff --git ") {
			break
		}

		pos++
	}

	// If binary file, we're done with this file
	if fd.IsBinary {
		return fd, pos, nil
	}

	// Parse --- and +++ lines
	if pos < len(lines) && strings.HasPrefix(lines[pos], "--- ") {
		path := strings.TrimPrefix(lines[pos], "--- ")
		if path == "/dev/null" {
			fd.IsNew = true
		} else {
			// Remove a/ prefix if present
			fd.OldPath = strings.TrimPrefix(path, "a/")
		}
		pos++
	}

	if pos < len(lines) && strings.HasPrefix(lines[pos], "+++ ") {
		path := strings.TrimPrefix(lines[pos], "+++ ")
		if path == "/dev/null" {
			fd.IsDeleted = true
			fd.Name = fd.OldPath
		} else {
			// Remove b/ prefix if present
			fd.NewPath = strings.TrimPrefix(path, "b/")
			fd.Name = fd.NewPath
		}
		pos++
	}

	// Parse hunks
	var hunks []hunk
	for pos < len(lines) {
		line := lines[pos]

		// Check for next file diff
		if strings.HasPrefix(line, "diff --git ") {
			break
		}

		// Check for hunk header
		if strings.HasPrefix(line, "@@ ") {
			h, newPos, err := parseHunk(lines, pos)
			if err != nil {
				return nil, newPos, err
			}
			hunks = append(hunks, h)
			pos = newPos
			continue
		}

		pos++
	}

	// Align hunks into left/right lines
	fd.LeftLines, fd.RightLines, fd.AddCount, fd.DelCount = alignHunks(hunks)

	return fd, pos, nil
}

// parseHunk parses a single hunk starting at pos (which should be at @@ line)
func parseHunk(lines []string, pos int) (hunk, int, error) {
	h := hunk{}

	// Parse header
	matches := hunkHeaderRE.FindStringSubmatch(lines[pos])
	if matches == nil {
		return h, pos + 1, &ParseError{Line: pos + 1, Message: "invalid hunk header"}
	}

	h.oldStart, _ = strconv.Atoi(matches[1])
	if matches[2] != "" {
		h.oldCount, _ = strconv.Atoi(matches[2])
	} else {
		h.oldCount = 1
	}
	h.newStart, _ = strconv.Atoi(matches[3])
	if matches[4] != "" {
		h.newCount, _ = strconv.Atoi(matches[4])
	} else {
		h.newCount = 1
	}

	pos++

	// Parse hunk lines
	for pos < len(lines) {
		line := lines[pos]

		// Check for next hunk or file
		if strings.HasPrefix(line, "@@ ") || strings.HasPrefix(line, "diff --git ") {
			break
		}

		// Handle "\ No newline at end of file"
		if strings.HasPrefix(line, "\\ ") {
			pos++
			continue
		}

		// Empty line in diff usually means context (space prefix got trimmed)
		if line == "" {
			// Could be end of diff or empty context line
			// Check if next line continues the diff
			if pos+1 < len(lines) {
				nextLine := lines[pos+1]
				if strings.HasPrefix(nextLine, " ") || strings.HasPrefix(nextLine, "+") ||
					strings.HasPrefix(nextLine, "-") || strings.HasPrefix(nextLine, "@@ ") ||
					strings.HasPrefix(nextLine, "diff --git ") || strings.HasPrefix(nextLine, "\\ ") {
					// Treat as empty context line
					h.lines = append(h.lines, diff.Line{Type: diff.Context, Content: ""})
					pos++
					continue
				}
			}
			// End of this file's diff
			break
		}

		// Parse line based on prefix
		if len(line) > 0 {
			switch line[0] {
			case ' ':
				h.lines = append(h.lines, diff.Line{Type: diff.Context, Content: line[1:]})
			case '+':
				h.lines = append(h.lines, diff.Line{Type: diff.Add, Content: line[1:]})
			case '-':
				h.lines = append(h.lines, diff.Line{Type: diff.Delete, Content: line[1:]})
			default:
				// Treat as context if no recognized prefix (shouldn't happen normally)
				h.lines = append(h.lines, diff.Line{Type: diff.Context, Content: line})
			}
		}

		pos++
	}

	return h, pos, nil
}
