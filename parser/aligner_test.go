package parser

import (
	"testing"

	"diff-tui/diff"
)

func TestAlignHunk_ContextOnly(t *testing.T) {
	h := hunk{
		lines: []diff.Line{
			{Type: diff.Context, Content: "line1"},
			{Type: diff.Context, Content: "line2"},
		},
	}

	left, right, adds, dels := alignHunk(h)

	if adds != 0 || dels != 0 {
		t.Errorf("expected 0 adds/dels, got %d/%d", adds, dels)
	}

	if len(left) != 2 || len(right) != 2 {
		t.Fatalf("expected 2 lines each side, got %d/%d", len(left), len(right))
	}

	// Context lines should appear on both sides
	for i := 0; i < 2; i++ {
		if left[i].Type != diff.Context || right[i].Type != diff.Context {
			t.Errorf("line %d: expected Context on both sides", i)
		}
		if left[i].Content != right[i].Content {
			t.Errorf("line %d: content mismatch", i)
		}
	}
}

func TestAlignHunk_DeleteOnly(t *testing.T) {
	h := hunk{
		lines: []diff.Line{
			{Type: diff.Delete, Content: "deleted1"},
			{Type: diff.Delete, Content: "deleted2"},
		},
	}

	left, right, adds, dels := alignHunk(h)

	if adds != 0 {
		t.Errorf("expected 0 adds, got %d", adds)
	}
	if dels != 2 {
		t.Errorf("expected 2 dels, got %d", dels)
	}

	if len(left) != 2 || len(right) != 2 {
		t.Fatalf("expected 2 lines each side, got %d/%d", len(left), len(right))
	}

	// Left should have deletes, right should have placeholders
	for i := 0; i < 2; i++ {
		if left[i].Type != diff.Delete {
			t.Errorf("left line %d: expected Delete, got %v", i, left[i].Type)
		}
		if right[i].Type != diff.Placeholder {
			t.Errorf("right line %d: expected Placeholder, got %v", i, right[i].Type)
		}
	}
}

func TestAlignHunk_AddOnly(t *testing.T) {
	h := hunk{
		lines: []diff.Line{
			{Type: diff.Add, Content: "added1"},
			{Type: diff.Add, Content: "added2"},
		},
	}

	left, right, adds, dels := alignHunk(h)

	if adds != 2 {
		t.Errorf("expected 2 adds, got %d", adds)
	}
	if dels != 0 {
		t.Errorf("expected 0 dels, got %d", dels)
	}

	if len(left) != 2 || len(right) != 2 {
		t.Fatalf("expected 2 lines each side, got %d/%d", len(left), len(right))
	}

	// Left should have placeholders, right should have adds
	for i := 0; i < 2; i++ {
		if left[i].Type != diff.Placeholder {
			t.Errorf("left line %d: expected Placeholder, got %v", i, left[i].Type)
		}
		if right[i].Type != diff.Add {
			t.Errorf("right line %d: expected Add, got %v", i, right[i].Type)
		}
	}
}

func TestAlignHunk_Modification(t *testing.T) {
	// Delete followed by Add = modification
	h := hunk{
		lines: []diff.Line{
			{Type: diff.Delete, Content: "old"},
			{Type: diff.Add, Content: "new"},
		},
	}

	left, right, adds, dels := alignHunk(h)

	if adds != 1 || dels != 1 {
		t.Errorf("expected 1 add/1 del, got %d/%d", adds, dels)
	}

	if len(left) != 1 || len(right) != 1 {
		t.Fatalf("expected 1 line each side, got %d/%d", len(left), len(right))
	}

	// Should be side by side: delete on left, add on right
	if left[0].Type != diff.Delete {
		t.Errorf("left: expected Delete, got %v", left[0].Type)
	}
	if right[0].Type != diff.Add {
		t.Errorf("right: expected Add, got %v", right[0].Type)
	}
}

func TestAlignHunk_MoreDeletesThanAdds(t *testing.T) {
	h := hunk{
		lines: []diff.Line{
			{Type: diff.Delete, Content: "del1"},
			{Type: diff.Delete, Content: "del2"},
			{Type: diff.Delete, Content: "del3"},
			{Type: diff.Add, Content: "add1"},
		},
	}

	left, right, adds, dels := alignHunk(h)

	if adds != 1 || dels != 3 {
		t.Errorf("expected 1 add/3 dels, got %d/%d", adds, dels)
	}

	if len(left) != 3 || len(right) != 3 {
		t.Fatalf("expected 3 lines each side, got %d/%d", len(left), len(right))
	}

	// First row: del1 | add1
	if left[0].Type != diff.Delete || right[0].Type != diff.Add {
		t.Error("row 0: expected Delete | Add")
	}

	// Second row: del2 | placeholder
	if left[1].Type != diff.Delete || right[1].Type != diff.Placeholder {
		t.Error("row 1: expected Delete | Placeholder")
	}

	// Third row: del3 | placeholder
	if left[2].Type != diff.Delete || right[2].Type != diff.Placeholder {
		t.Error("row 2: expected Delete | Placeholder")
	}
}

func TestAlignHunk_MoreAddsThanDeletes(t *testing.T) {
	h := hunk{
		lines: []diff.Line{
			{Type: diff.Delete, Content: "del1"},
			{Type: diff.Add, Content: "add1"},
			{Type: diff.Add, Content: "add2"},
			{Type: diff.Add, Content: "add3"},
		},
	}

	left, right, adds, dels := alignHunk(h)

	if adds != 3 || dels != 1 {
		t.Errorf("expected 3 adds/1 del, got %d/%d", adds, dels)
	}

	if len(left) != 3 || len(right) != 3 {
		t.Fatalf("expected 3 lines each side, got %d/%d", len(left), len(right))
	}

	// First row: del1 | add1
	if left[0].Type != diff.Delete || right[0].Type != diff.Add {
		t.Error("row 0: expected Delete | Add")
	}

	// Second row: placeholder | add2
	if left[1].Type != diff.Placeholder || right[1].Type != diff.Add {
		t.Error("row 1: expected Placeholder | Add")
	}

	// Third row: placeholder | add3
	if left[2].Type != diff.Placeholder || right[2].Type != diff.Add {
		t.Error("row 2: expected Placeholder | Add")
	}
}

func TestAlignHunk_MixedWithContext(t *testing.T) {
	h := hunk{
		lines: []diff.Line{
			{Type: diff.Context, Content: "ctx1"},
			{Type: diff.Delete, Content: "del"},
			{Type: diff.Add, Content: "add"},
			{Type: diff.Context, Content: "ctx2"},
		},
	}

	left, right, adds, dels := alignHunk(h)

	if adds != 1 || dels != 1 {
		t.Errorf("expected 1/1, got %d/%d", adds, dels)
	}

	if len(left) != 3 || len(right) != 3 {
		t.Fatalf("expected 3 lines each side, got %d/%d", len(left), len(right))
	}

	// Row 0: ctx1 | ctx1
	if left[0].Type != diff.Context || right[0].Type != diff.Context {
		t.Error("row 0: expected Context | Context")
	}

	// Row 1: del | add (modification)
	if left[1].Type != diff.Delete || right[1].Type != diff.Add {
		t.Error("row 1: expected Delete | Add")
	}

	// Row 2: ctx2 | ctx2
	if left[2].Type != diff.Context || right[2].Type != diff.Context {
		t.Error("row 2: expected Context | Context")
	}
}

func TestAlignHunks_MultipleHunks(t *testing.T) {
	hunks := []hunk{
		{
			lines: []diff.Line{
				{Type: diff.Context, Content: "ctx1"},
				{Type: diff.Delete, Content: "del1"},
			},
		},
		{
			lines: []diff.Line{
				{Type: diff.Add, Content: "add2"},
				{Type: diff.Context, Content: "ctx2"},
			},
		},
	}

	left, right, adds, dels := alignHunks(hunks)

	if adds != 1 || dels != 1 {
		t.Errorf("expected 1/1, got %d/%d", adds, dels)
	}

	// Should have combined lines from both hunks
	if len(left) != 4 || len(right) != 4 {
		t.Fatalf("expected 4 lines each side, got %d/%d", len(left), len(right))
	}
}
