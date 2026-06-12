package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func keyMsg(t tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: t}
}

func runeMsg(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func applyKeys(m Model, msgs ...tea.KeyMsg) Model {
	for _, msg := range msgs {
		newModel, _ := m.Update(msg)
		m = newModel.(Model)
	}
	return m
}

func TestDefaultKeyMapUndoRedo(t *testing.T) {
	km := DefaultKeyMap()

	if !matches(keyMsg(tea.KeyCtrlZ), km.Undo) {
		t.Error("Ctrl+Z should match Undo")
	}
	if !matches(keyMsg(tea.KeyCtrlY), km.Redo) {
		t.Error("Ctrl+Y should match Redo")
	}
	if !matches(keyMsg(tea.KeyCtrlR), km.Redo) {
		t.Error("Ctrl+R should match Redo")
	}
	if matches(keyMsg(tea.KeyCtrlZ), km.Redo) {
		t.Error("Ctrl+Z should not match Redo")
	}
}

func TestUndoEnter(t *testing.T) {
	m := New()
	initialCount := len(m.nodes)

	m = applyKeys(m, keyMsg(tea.KeyEnter))
	if len(m.nodes) != initialCount+1 {
		t.Fatalf("expected %d nodes after enter, got %d", initialCount+1, len(m.nodes))
	}

	m = applyKeys(m, keyMsg(tea.KeyCtrlZ))
	if len(m.nodes) != initialCount {
		t.Errorf("expected %d nodes after undo, got %d", initialCount, len(m.nodes))
	}
}

func TestUndoRestoresCursor(t *testing.T) {
	m := New()

	m = applyKeys(m, keyMsg(tea.KeyEnter))
	if m.cursor != 1 {
		t.Fatalf("expected cursor at 1 after enter, got %d", m.cursor)
	}

	m = applyKeys(m, keyMsg(tea.KeyCtrlZ))
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0 after undo, got %d", m.cursor)
	}
}

func TestUndoDelete(t *testing.T) {
	m := New()
	m = applyKeys(m, runeMsg('a'), keyMsg(tea.KeyEnter), runeMsg('b'))
	countBefore := len(m.nodes)

	m = applyKeys(m, keyMsg(tea.KeyCtrlD))
	if len(m.nodes) != countBefore-1 {
		t.Fatalf("expected %d nodes after delete, got %d", countBefore-1, len(m.nodes))
	}

	m = applyKeys(m, keyMsg(tea.KeyCtrlZ))
	if len(m.nodes) != countBefore {
		t.Errorf("expected %d nodes after undo, got %d", countBefore, len(m.nodes))
	}
	if m.nodes[1].Text != "b" {
		t.Errorf("expected restored node text 'b', got %q", m.nodes[1].Text)
	}
}

func TestUndoIndent(t *testing.T) {
	m := New()
	m = applyKeys(m, runeMsg('a'), keyMsg(tea.KeyEnter), runeMsg('b'))

	m = applyKeys(m, keyMsg(tea.KeyTab))
	if m.currentNode().Depth() != 2 {
		t.Fatalf("expected depth 2 after indent, got %d", m.currentNode().Depth())
	}

	m = applyKeys(m, keyMsg(tea.KeyCtrlZ))
	if m.currentNode().Depth() != 1 {
		t.Errorf("expected depth 1 after undo, got %d", m.currentNode().Depth())
	}
}

func TestUndoUnindent(t *testing.T) {
	m := New()
	m = applyKeys(m, runeMsg('a'), keyMsg(tea.KeyEnter), runeMsg('b'), keyMsg(tea.KeyTab))
	if m.currentNode().Depth() != 2 {
		t.Fatalf("expected depth 2 after indent, got %d", m.currentNode().Depth())
	}

	m = applyKeys(m, keyMsg(tea.KeyShiftTab))
	if m.currentNode().Depth() != 1 {
		t.Fatalf("expected depth 1 after unindent, got %d", m.currentNode().Depth())
	}

	m = applyKeys(m, keyMsg(tea.KeyCtrlZ))
	if m.currentNode().Depth() != 2 {
		t.Errorf("expected depth 2 after undo, got %d", m.currentNode().Depth())
	}
}

func TestRedo(t *testing.T) {
	m := New()
	initialCount := len(m.nodes)

	m = applyKeys(m, keyMsg(tea.KeyEnter), keyMsg(tea.KeyCtrlZ))
	if len(m.nodes) != initialCount {
		t.Fatalf("expected %d nodes after undo, got %d", initialCount, len(m.nodes))
	}

	m = applyKeys(m, keyMsg(tea.KeyCtrlY))
	if len(m.nodes) != initialCount+1 {
		t.Errorf("expected %d nodes after redo, got %d", initialCount+1, len(m.nodes))
	}
}

func TestRedoClearedByNewEdit(t *testing.T) {
	m := New()

	m = applyKeys(m, keyMsg(tea.KeyEnter), keyMsg(tea.KeyCtrlZ))
	// A new edit after undo must clear the redo stack
	m = applyKeys(m, runeMsg('x'))
	countAfterEdit := len(m.nodes)

	m = applyKeys(m, keyMsg(tea.KeyCtrlY))
	if len(m.nodes) != countAfterEdit {
		t.Errorf("redo after new edit should do nothing, got %d nodes (was %d)", len(m.nodes), countAfterEdit)
	}
	if m.currentNode().Text != "x" {
		t.Errorf("expected text 'x' preserved, got %q", m.currentNode().Text)
	}
}

func TestUndoTextEditCoalescing(t *testing.T) {
	m := New()

	// Typing consecutive characters into the same node is one undo step
	m = applyKeys(m, runeMsg('a'), runeMsg('b'), runeMsg('c'))
	if m.currentNode().Text != "abc" {
		t.Fatalf("expected text 'abc', got %q", m.currentNode().Text)
	}

	m = applyKeys(m, keyMsg(tea.KeyCtrlZ))
	if m.currentNode().Text != "" {
		t.Errorf("expected empty text after single undo, got %q", m.currentNode().Text)
	}
}

func TestUndoTextEditSeparateNodes(t *testing.T) {
	m := New()

	// Edits on different nodes are separate undo steps
	m = applyKeys(m, runeMsg('a'), keyMsg(tea.KeyEnter), runeMsg('b'))

	m = applyKeys(m, keyMsg(tea.KeyCtrlZ))
	if m.currentNode().Text != "" {
		t.Errorf("expected second node text cleared, got %q", m.currentNode().Text)
	}
	if m.nodes[0].Text != "a" {
		t.Errorf("expected first node text 'a' intact, got %q", m.nodes[0].Text)
	}
}

func TestUndoNothingToUndo(t *testing.T) {
	m := New()

	m = applyKeys(m, keyMsg(tea.KeyCtrlZ))
	if m.message != "Nothing to undo" {
		t.Errorf("expected 'Nothing to undo' message, got %q", m.message)
	}
}

func TestRedoNothingToRedo(t *testing.T) {
	m := New()

	m = applyKeys(m, keyMsg(tea.KeyCtrlY))
	if m.message != "Nothing to redo" {
		t.Errorf("expected 'Nothing to redo' message, got %q", m.message)
	}
}

func TestUndoSyncsTextInput(t *testing.T) {
	m := New()

	m = applyKeys(m, runeMsg('a'), keyMsg(tea.KeyEnter), runeMsg('b'), keyMsg(tea.KeyCtrlZ))
	// After undoing the 'b' edit, the text input must show the restored text
	if m.textInput.Value() != m.currentNode().Text {
		t.Errorf("text input %q out of sync with node text %q", m.textInput.Value(), m.currentNode().Text)
	}
}

func TestUndoHistoryLimit(t *testing.T) {
	m := New()

	// Push more snapshots than the history cap
	for i := 0; i < maxHistory+10; i++ {
		m = applyKeys(m, keyMsg(tea.KeyEnter))
	}

	if len(m.history.undo) > maxHistory {
		t.Errorf("undo stack should be capped at %d, got %d", maxHistory, len(m.history.undo))
	}
}
