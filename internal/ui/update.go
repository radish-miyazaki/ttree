package ui

import (
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/radish-miyazaki/ttree/internal/tree"
)

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Handle text input updates
	if m.mode == ModeEdit {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		// Sync text to node in real-time
		if node := m.currentNode(); node != nil {
			node.Text = m.textInput.Value()
		}
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.copied = false
	m.message = ""

	// Handle quit
	if matches(msg, m.keys.Quit) {
		m.saveCurrentEdit()
		return m, tea.Quit
	}

	// Handle copy
	if matches(msg, m.keys.Copy) {
		m.saveCurrentEdit()
		output := m.renderer.Render(m.tree)
		if err := clipboard.WriteAll(output); err == nil {
			m.copied = true
			m.message = "Copied to clipboard!"
		} else {
			m.message = "Failed to copy: " + err.Error()
		}
		return m, nil
	}

	// Navigation
	if matches(msg, m.keys.Up) {
		m.moveCursor(-1)
		return m, nil
	}
	if matches(msg, m.keys.Down) {
		m.moveCursor(1)
		return m, nil
	}

	// Indent / Unindent
	if matches(msg, m.keys.Indent) {
		m.saveCurrentEdit()
		if node := m.currentNode(); node != nil {
			m.tree.Indent(node)
			m.refreshNodes()
			m.focusNode(node)
		}
		return m, nil
	}
	if matches(msg, m.keys.Unindent) {
		m.saveCurrentEdit()
		if node := m.currentNode(); node != nil {
			m.tree.Unindent(node)
			m.refreshNodes()
			m.focusNode(node)
		}
		return m, nil
	}

	// Enter - create new sibling
	if matches(msg, m.keys.Enter) {
		m.saveCurrentEdit()
		if node := m.currentNode(); node != nil {
			newNode := tree.NewNode("")
			m.tree.InsertAfter(node, newNode)
			m.refreshNodes()
			m.focusNode(newNode)
		}
		return m, nil
	}

	// Delete
	if matches(msg, m.keys.Delete) {
		m.saveCurrentEdit()
		if node := m.currentNode(); node != nil {
			nextFocus := m.tree.Delete(node)
			m.refreshNodes()
			if nextFocus != nil {
				m.focusNode(nextFocus)
			}
		}
		return m, nil
	}

	// Pass to text input
	if m.mode == ModeEdit {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		if node := m.currentNode(); node != nil {
			node.Text = m.textInput.Value()
		}
		return m, cmd
	}

	return m, nil
}
