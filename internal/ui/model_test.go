package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	m := New()

	if m.tree == nil {
		t.Error("expected non-nil tree")
	}
	if m.renderer == nil {
		t.Error("expected non-nil renderer")
	}
	if len(m.nodes) == 0 {
		t.Error("expected at least one node")
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", m.cursor)
	}
	if m.mode != ModeEdit {
		t.Error("expected edit mode by default")
	}
}

func TestCurrentNode(t *testing.T) {
	m := New()
	node := m.currentNode()

	if node == nil {
		t.Error("expected non-nil current node")
	}
}

func TestCurrentNodeOutOfBounds(t *testing.T) {
	m := New()
	m.cursor = 999

	node := m.currentNode()
	if node != nil {
		t.Error("expected nil for out of bounds cursor")
	}

	m.cursor = -1
	node = m.currentNode()
	if node != nil {
		t.Error("expected nil for negative cursor")
	}
}

func TestRefreshNodes(t *testing.T) {
	m := New()
	initialCount := len(m.nodes)

	// Add a node
	m.tree.Root.Children[0].Text = "test"
	m.tree.InsertAfter(m.nodes[0], m.tree.Root.Children[0])
	m.refreshNodes()

	if len(m.nodes) <= initialCount {
		t.Error("expected more nodes after insert")
	}
}

func TestRefreshNodesCursorBounds(t *testing.T) {
	m := New()
	m.cursor = 100

	m.refreshNodes()

	if m.cursor >= len(m.nodes) {
		t.Error("cursor should be bounded to valid range")
	}
}

func TestMoveCursor(t *testing.T) {
	m := New()

	// Add more nodes
	m.currentNode().Text = "node1"
	m.tree.InsertAfter(m.currentNode(), m.tree.Root.Children[0])
	m.refreshNodes()

	if len(m.nodes) < 2 {
		t.Skip("need at least 2 nodes for this test")
	}

	m.moveCursor(1)
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1, got %d", m.cursor)
	}

	m.moveCursor(-1)
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", m.cursor)
	}
}

func TestMoveCursorBounds(t *testing.T) {
	m := New()

	// Try to move up from top
	m.moveCursor(-1)
	if m.cursor != 0 {
		t.Error("cursor should not go negative")
	}

	// Try to move past end
	m.moveCursor(100)
	if m.cursor >= len(m.nodes) {
		t.Error("cursor should not exceed node count")
	}
}

func TestSaveCurrentEdit(t *testing.T) {
	m := New()
	m.textInput.SetValue("new text")

	m.saveCurrentEdit()

	if m.currentNode().Text != "new text" {
		t.Errorf("expected 'new text', got %q", m.currentNode().Text)
	}
}

func TestInit(t *testing.T) {
	m := New()
	cmd := m.Init()

	if cmd == nil {
		t.Error("expected non-nil init command (textinput.Blink)")
	}
}

func TestUpdateWindowSize(t *testing.T) {
	m := New()
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}

	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.width != 100 {
		t.Errorf("expected width 100, got %d", updated.width)
	}
	if updated.height != 50 {
		t.Errorf("expected height 50, got %d", updated.height)
	}
}

func TestUpdateQuit(t *testing.T) {
	m := New()

	// Test ctrl+q
	msg := tea.KeyMsg{Type: tea.KeyCtrlQ}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestUpdateNavigation(t *testing.T) {
	m := New()

	// Add another node first
	m.currentNode().Text = "node1"
	m.tree.InsertAfter(m.currentNode(), m.tree.Root.Children[0])
	m.refreshNodes()

	if len(m.nodes) < 2 {
		t.Skip("need at least 2 nodes")
	}

	// Test down navigation
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if updated.cursor != 1 {
		t.Errorf("expected cursor at 1 after down, got %d", updated.cursor)
	}

	// Test up navigation
	msg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = updated.Update(msg)
	updated = newModel.(Model)

	if updated.cursor != 0 {
		t.Errorf("expected cursor at 0 after up, got %d", updated.cursor)
	}
}

func TestUpdateEnter(t *testing.T) {
	m := New()
	initialCount := len(m.nodes)

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if len(updated.nodes) != initialCount+1 {
		t.Errorf("expected %d nodes after enter, got %d", initialCount+1, len(updated.nodes))
	}
}

func TestUpdateIndent(t *testing.T) {
	m := New()

	// Add another node using Enter key
	m.currentNode().Text = "node1"
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(enterMsg)
	m = newModel.(Model)

	// Now cursor should be on the new (second) node
	if len(m.nodes) < 2 {
		t.Skip("need at least 2 nodes")
	}

	// Cursor should already be on the second node after Enter
	initialDepth := m.currentNode().Depth()

	msg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ = m.Update(msg)
	updated := newModel.(Model)

	newDepth := updated.currentNode().Depth()
	if newDepth != initialDepth+1 {
		t.Errorf("expected depth %d after indent, got %d", initialDepth+1, newDepth)
	}
}

func TestUpdateUnindent(t *testing.T) {
	m := New()

	// Setup: create second node using Enter
	m.currentNode().Text = "parent"
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := m.Update(enterMsg)
	m = newModel.(Model)

	// Indent the current (second) node
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ = m.Update(tabMsg)
	m = newModel.(Model)

	initialDepth := m.currentNode().Depth()
	if initialDepth < 2 {
		t.Skipf("node not indented properly, depth=%d", initialDepth)
	}

	// Now unindent
	msg := tea.KeyMsg{Type: tea.KeyShiftTab}
	newModel, _ = m.Update(msg)
	updated := newModel.(Model)

	if updated.currentNode().Depth() >= initialDepth {
		t.Errorf("depth should decrease after unindent, was %d, now %d", initialDepth, updated.currentNode().Depth())
	}
}

func TestUpdateDelete(t *testing.T) {
	m := New()

	// Add nodes first
	m.currentNode().Text = "node1"
	m.tree.InsertAfter(m.currentNode(), m.tree.Root.Children[0])
	m.refreshNodes()

	initialCount := len(m.nodes)
	if initialCount < 2 {
		t.Skip("need at least 2 nodes")
	}

	msg := tea.KeyMsg{Type: tea.KeyCtrlD}
	newModel, _ := m.Update(msg)
	updated := newModel.(Model)

	if len(updated.nodes) != initialCount-1 {
		t.Errorf("expected %d nodes after delete, got %d", initialCount-1, len(updated.nodes))
	}
}

func TestViewNotEmpty(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	view := m.View()

	if view == "" {
		t.Error("expected non-empty view")
	}
	if view == "Loading..." {
		t.Error("view should not be loading with set dimensions")
	}
}

func TestViewLoading(t *testing.T) {
	m := New()
	m.width = 0

	view := m.View()

	if view != "Loading..." {
		t.Errorf("expected 'Loading...' when width is 0, got %q", view)
	}
}
