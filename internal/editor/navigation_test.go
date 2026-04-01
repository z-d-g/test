package editor

import "testing"

func newTestNav(content string) *Navigation {
	buf := NewGapBuffer([]byte(content))
	return NewNavigation(buf, 0)
}

func TestNewNavigation(t *testing.T) {
	nav := newTestNav("hello")
	if nav.Cursor() != 0 {
		t.Errorf("initial cursor = %d, want 0", nav.Cursor())
	}
}

func TestSetCursor(t *testing.T) {
	nav := newTestNav("hello world")

	nav.SetCursor(5)
	if nav.Cursor() != 5 {
		t.Errorf("cursor = %d, want 5", nav.Cursor())
	}

	nav.SetCursor(-1)
	if nav.Cursor() != 0 {
		t.Errorf("negative cursor = %d, want 0", nav.Cursor())
	}

	nav.SetCursor(100)
	if nav.Cursor() != nav.buf.Len() {
		t.Errorf("overshoot cursor = %d, want %d", nav.Cursor(), nav.buf.Len())
	}
}

func TestMoveRight(t *testing.T) {
	nav := newTestNav("abc")

	nav.MoveRight()
	if nav.Cursor() != 1 {
		t.Errorf("cursor = %d, want 1", nav.Cursor())
	}

	nav.MoveRight()
	nav.MoveRight()
	if nav.Cursor() != 3 {
		t.Errorf("cursor = %d, want 3", nav.Cursor())
	}

	// At end, should not move
	nav.MoveRight()
	if nav.Cursor() != 3 {
		t.Errorf("cursor at end = %d, want 3", nav.Cursor())
	}
}

func TestMoveLeft(t *testing.T) {
	nav := newTestNav("abc")
	nav.SetCursor(3)

	nav.MoveLeft()
	if nav.Cursor() != 2 {
		t.Errorf("cursor = %d, want 2", nav.Cursor())
	}

	nav.MoveLeft()
	nav.MoveLeft()
	if nav.Cursor() != 0 {
		t.Errorf("cursor = %d, want 0", nav.Cursor())
	}

	// At start, should not move
	nav.MoveLeft()
	if nav.Cursor() != 0 {
		t.Errorf("cursor at start = %d, want 0", nav.Cursor())
	}
}

func TestMoveUp(t *testing.T) {
	nav := newTestNav("line1\nline2\nline3")
	nav.SetCursor(10) // on line2

	nav.MoveUp()
	row, _ := nav.buf.CursorToRowCol(nav.Cursor())
	if row != 0 {
		t.Errorf("row after up = %d, want 0", row)
	}

	// At top, should not move
	nav.MoveUp()
	row, _ = nav.buf.CursorToRowCol(nav.Cursor())
	if row != 0 {
		t.Errorf("row at top = %d, want 0", row)
	}
}

func TestMoveDown(t *testing.T) {
	nav := newTestNav("line1\nline2\nline3")
	nav.SetCursor(2) // on line1

	nav.MoveDown()
	row, _ := nav.buf.CursorToRowCol(nav.Cursor())
	if row != 1 {
		t.Errorf("row after down = %d, want 1", row)
	}

	nav.MoveDown()
	row, _ = nav.buf.CursorToRowCol(nav.Cursor())
	if row != 2 {
		t.Errorf("row after second down = %d, want 2", row)
	}

	// At bottom, should not move
	nav.MoveDown()
	row, _ = nav.buf.CursorToRowCol(nav.Cursor())
	if row != 2 {
		t.Errorf("row at bottom = %d, want 2", row)
	}
}

func TestMoveHome(t *testing.T) {
	nav := newTestNav("hello\nworld")
	nav.SetCursor(8) // middle of line2

	nav.MoveHome()
	if nav.Cursor() != 6 {
		t.Errorf("cursor = %d, want 6 (start of line 2)", nav.Cursor())
	}
}

func TestMoveEnd(t *testing.T) {
	nav := newTestNav("hello\nworld")
	nav.SetCursor(6) // start of line2

	nav.MoveEnd()
	if nav.Cursor() != 11 {
		t.Errorf("cursor = %d, want 11 (end of doc)", nav.Cursor())
	}
}

func TestMoveDocStart(t *testing.T) {
	nav := newTestNav("hello\nworld")
	nav.SetCursor(5)

	nav.MoveDocStart()
	if nav.Cursor() != 0 {
		t.Errorf("cursor = %d, want 0", nav.Cursor())
	}
}

func TestMoveDocEnd(t *testing.T) {
	nav := newTestNav("hello\nworld")

	nav.MoveDocEnd()
	if nav.Cursor() != 11 {
		t.Errorf("cursor = %d, want 11", nav.Cursor())
	}
}

func TestMoveWordLeft(t *testing.T) {
	nav := newTestNav("hello world")
	nav.SetCursor(11)

	nav.MoveWordLeft()
	if nav.Cursor() != 6 {
		t.Errorf("cursor = %d, want 6", nav.Cursor())
	}

	nav.MoveWordLeft()
	if nav.Cursor() != 0 {
		t.Errorf("cursor = %d, want 0", nav.Cursor())
	}

	// At start
	nav.MoveWordLeft()
	if nav.Cursor() != 0 {
		t.Errorf("cursor at start = %d, want 0", nav.Cursor())
	}
}

func TestMoveWordRight(t *testing.T) {
	nav := newTestNav("hello world")

	nav.MoveWordRight()
	if nav.Cursor() != 6 {
		t.Errorf("cursor = %d, want 6", nav.Cursor())
	}

	nav.MoveWordRight()
	if nav.Cursor() != 11 {
		t.Errorf("cursor = %d, want 11", nav.Cursor())
	}

	// At end
	nav.MoveWordRight()
	if nav.Cursor() != 11 {
		t.Errorf("cursor at end = %d, want 11", nav.Cursor())
	}
}

func TestPageUp(t *testing.T) {
	nav := newTestNav("l1\nl2\nl3\nl4\nl5\nl6\nl7\nl8\nl9\nl10")
	nav.SetCursor(nav.buf.Len()) // end

	nav.PageUp(5)
	row, _ := nav.buf.CursorToRowCol(nav.Cursor())
	if row != 4 {
		t.Errorf("row after page up = %d, want 4", row)
	}

	// Invalid viewport
	nav.PageUp(0)
	nav.PageUp(-1)
}

func TestPageDown(t *testing.T) {
	nav := newTestNav("l1\nl2\nl3\nl4\nl5\nl6\nl7\nl8\nl9\nl10")

	nav.PageDown(5)
	row, _ := nav.buf.CursorToRowCol(nav.Cursor())
	if row != 5 {
		t.Errorf("row after page down = %d, want 5", row)
	}

	// Invalid viewport
	nav.PageDown(0)
	nav.PageDown(-1)
}
