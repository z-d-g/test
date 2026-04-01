package editor

import (
	"testing"
)

func TestSelectionVisualization(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3\nLine 4"
	editor := NewEditor(content, testRenderer())

	editor.selection.Activate(0)
	editor.selection.Extend(4)

	if !editor.selection.IsActive() {
		t.Error("Selection should be active")
	}

	if editor.selection.Start() != 0 {
		t.Errorf("Selection start should be 0, got %d", editor.selection.Start())
	}

	if editor.selection.End() != 4 {
		t.Errorf("Selection end should be 4, got %d", editor.selection.End())
	}

	editor.selection.Activate(0)
	editor.selection.Extend(12)

	if editor.selection.Start() != 0 {
		t.Errorf("Multi-line selection start should be 0, got %d", editor.selection.Start())
	}

	if editor.selection.End() != 12 {
		t.Errorf("Multi-line selection end should be 12, got %d", editor.selection.End())
	}

	editor.selection.Clear()

	if editor.selection.IsActive() {
		t.Error("Selection should not be active after Clear()")
	}
}

func TestSelectionCut(t *testing.T) {
	content := "HelloWorld"
	editor := NewEditor(content, testRenderer())

	editor.selection.Activate(5)
	editor.selection.Extend(10)

	cutContent := editor.selection.Cut(editor.buf)

	if string(cutContent) != "World" {
		t.Errorf("Cut content should be 'World', got '%s'", string(cutContent))
	}

	if editor.buf.Len() != 5 {
		t.Errorf("Buffer length should be 5 after cut, got %d", editor.buf.Len())
	}

	if string(editor.buf.Contents()) != "Hello" {
		t.Errorf("Buffer content should be 'Hello', got '%s'", string(editor.buf.Contents()))
	}

	if editor.selection.IsActive() {
		t.Error("Selection should not be active after Cut()")
	}
}

func TestSelectionGetSelectedText(t *testing.T) {
	content := "Hello World"
	editor := NewEditor(content, testRenderer())

	editor.selection.Activate(6)
	editor.selection.Extend(11)

	selected := editor.selection.GetSelectedText(editor.buf)

	if string(selected) != "World" {
		t.Errorf("Selected text should be 'World', got '%s'", string(selected))
	}

	if !editor.selection.IsActive() {
		t.Error("Selection should remain active after GetSelectedText()")
	}
}

func TestSelectionLength(t *testing.T) {
	sel := NewSelection()
	if sel.Length() != 0 {
		t.Errorf("inactive selection length = %d, want 0", sel.Length())
	}

	sel.Activate(0)
	sel.Extend(5)
	if sel.Length() != 5 {
		t.Errorf("selection length = %d, want 5", sel.Length())
	}

	// Reverse direction
	sel.Activate(10)
	sel.Extend(3)
	if sel.Length() != 7 {
		t.Errorf("reverse selection length = %d, want 7", sel.Length())
	}
}

func TestSelectionSelectAll(t *testing.T) {
	content := "Hello World"
	editor := NewEditor(content, testRenderer())

	editor.selection.SelectAll(editor.buf)

	if !editor.selection.IsActive() {
		t.Error("selection should be active")
	}
	if editor.selection.Start() != 0 {
		t.Errorf("start = %d, want 0", editor.selection.Start())
	}
	if editor.selection.End() != len(content) {
		t.Errorf("end = %d, want %d", editor.selection.End(), len(content))
	}
}

func TestSelectionSelectWord(t *testing.T) {
	content := "hello world"
	editor := NewEditor(content, testRenderer())

	editor.selection.SelectWord(editor.buf, 2) // inside "hello"

	if editor.selection.Start() != 0 {
		t.Errorf("word start = %d, want 0", editor.selection.Start())
	}
	if editor.selection.End() != 5 {
		t.Errorf("word end = %d, want 5", editor.selection.End())
	}

	// Select second word
	editor.selection.SelectWord(editor.buf, 7) // inside "world"
	if editor.selection.Start() != 6 {
		t.Errorf("word start = %d, want 6", editor.selection.Start())
	}
	if editor.selection.End() != 11 {
		t.Errorf("word end = %d, want 11", editor.selection.End())
	}
}

func TestSelectionSelectWordEmpty(t *testing.T) {
	editor := NewEditor("", testRenderer())
	editor.selection.SelectWord(editor.buf, 0)
	if editor.selection.IsActive() {
		t.Error("selecting word in empty buffer should not activate")
	}
}

func TestSelectionSelectLine(t *testing.T) {
	content := "line1\nline2\nline3"
	editor := NewEditor(content, testRenderer())

	editor.selection.SelectLine(editor.buf, 2) // inside first line

	if editor.selection.Start() != 0 {
		t.Errorf("line start = %d, want 0", editor.selection.Start())
	}
	if editor.selection.End() != 5 {
		t.Errorf("line end = %d, want 5", editor.selection.End())
	}
}

func TestSelectionMoveCursor(t *testing.T) {
	sel := NewSelection()

	// First move activates
	sel.MoveCursor(5)
	if !sel.IsActive() {
		t.Error("MoveCursor should activate selection")
	}
	if sel.Length() != 0 {
		t.Errorf("length = %d, want 0", sel.Length())
	}

	// Second move extends
	sel.MoveCursor(10)
	if sel.End() != 10 {
		t.Errorf("end = %d, want 10", sel.End())
	}
}

func TestSelectionGetSelectedTextClamp(t *testing.T) {
	content := "short"
	editor := NewEditor(content, testRenderer())

	editor.selection.Activate(0)
	editor.selection.Extend(100) // beyond buffer

	selected := editor.selection.GetSelectedText(editor.buf)
	if string(selected) != "short" {
		t.Errorf("clamped text = %q, want %q", selected, "short")
	}
}

func TestSelectionCutInactive(t *testing.T) {
	content := "hello"
	editor := NewEditor(content, testRenderer())

	result := editor.selection.Cut(editor.buf)
	if result != nil {
		t.Error("cut on inactive selection should return nil")
	}
}

func TestSelectionWithUndo(t *testing.T) {
	content := "HelloWorld"
	editor := NewEditor(content, testRenderer())

	editor.selection.Activate(5)
	editor.selection.Extend(10)

	kb := editor.keyBindings
	kb.handleCut()

	if string(editor.buf.Contents()) != "Hello" {
		t.Errorf("After cut, buffer content should be 'Hello', got '%s'", string(editor.buf.Contents()))
	}

	editor.undo.Undo(editor.buf)

	if string(editor.buf.Contents()) != "HelloWorld" {
		t.Errorf("After undo, buffer content should be 'HelloWorld', got '%s'", string(editor.buf.Contents()))
	}
}
