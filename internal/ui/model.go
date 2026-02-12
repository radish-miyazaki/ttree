package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/radish-miyazaki/ttree/internal/render"
	"github.com/radish-miyazaki/ttree/internal/tree"
)

// Mode represents the current editing mode
type Mode int

const (
	ModeNormal Mode = iota
	ModeEdit
)

// Model represents the application state
type Model struct {
	tree       *tree.Tree
	renderer   *render.Renderer
	cursor     int           // Current cursor position in flattened list
	nodes      []*tree.Node  // Flattened visible nodes
	mode       Mode
	textInput  textinput.Model
	width      int
	height     int
	keys       KeyMap
	message    string        // Status message
	copied     bool          // Flash message for copy
}

// New creates a new model
func New() Model {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Width = 50

	t := tree.NewTree()
	nodes := t.FlattenVisible()

	m := Model{
		tree:      t,
		renderer:  render.NewRenderer(),
		cursor:    0,
		nodes:     nodes,
		mode:      ModeEdit,
		textInput: ti,
		keys:      DefaultKeyMap(),
		width:     80,
		height:    24,
	}

	// Start editing the first node
	if len(nodes) > 0 {
		m.textInput.SetValue(nodes[0].Text)
		m.textInput.Focus()
	}

	return m
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// currentNode returns the currently selected node
func (m *Model) currentNode() *tree.Node {
	if m.cursor >= 0 && m.cursor < len(m.nodes) {
		return m.nodes[m.cursor]
	}
	return nil
}

// refreshNodes updates the flattened node list
func (m *Model) refreshNodes() {
	m.nodes = m.tree.FlattenVisible()
	if m.cursor >= len(m.nodes) {
		m.cursor = len(m.nodes) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// syncTextInput syncs the text input with the current node
func (m *Model) syncTextInput() {
	if node := m.currentNode(); node != nil {
		m.textInput.SetValue(node.Text)
		m.textInput.CursorEnd()
	}
}

// saveCurrentEdit saves the current text input to the node
func (m *Model) saveCurrentEdit() {
	if node := m.currentNode(); node != nil {
		node.Text = m.textInput.Value()
	}
}

// moveCursor moves the cursor and handles edit sync
func (m *Model) moveCursor(delta int) {
	m.saveCurrentEdit()
	newPos := m.cursor + delta
	if newPos >= 0 && newPos < len(m.nodes) {
		m.cursor = newPos
		m.syncTextInput()
	}
}

// focusNode sets cursor to specific node
func (m *Model) focusNode(n *tree.Node) {
	for i, node := range m.nodes {
		if node.ID == n.ID {
			m.cursor = i
			m.syncTextInput()
			return
		}
	}
}
