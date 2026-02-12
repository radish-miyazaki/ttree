package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Editor pane styles
	editorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	// Preview pane styles
	previewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	// Selected line style
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230"))

	// Normal line style
	normalStyle = lipgloss.NewStyle()

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62")).
			MarginBottom(0)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Status message style
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120"))
)

// View implements tea.Model
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Calculate pane widths
	totalWidth := m.width
	editorWidth := totalWidth/2 - 2
	previewWidth := totalWidth - editorWidth - 4

	// Build editor view
	editorContent := m.buildEditorView(editorWidth)
	editorPane := editorStyle.Width(editorWidth).Height(m.height - 5).Render(editorContent)

	// Build preview view
	previewContent := m.buildPreviewView()
	previewPane := previewStyle.Width(previewWidth).Height(m.height - 5).Render(previewContent)

	// Combine panes
	content := lipgloss.JoinHorizontal(lipgloss.Top, editorPane, previewPane)

	// Build help line
	help := m.buildHelpLine()

	// Build status line
	status := ""
	if m.message != "" {
		status = statusStyle.Render(m.message)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(" ttree - Tree Editor"),
		content,
		help,
		status,
	)
}

func (m Model) buildEditorView(width int) string {
	var lines []string

	for i, node := range m.nodes {
		// Build indentation
		depth := node.Depth()
		indent := strings.Repeat("  ", depth-1)

		// Build line content
		var line string
		if i == m.cursor {
			// Current line with text input
			prefix := indent + "• "
			inputWidth := width - len(prefix) - 2
			if inputWidth < 10 {
				inputWidth = 10
			}
			m.textInput.Width = inputWidth
			line = prefix + m.textInput.View()
		} else {
			// Regular line
			text := node.Text
			if text == "" {
				text = " "
			}
			line = indent + "• " + text
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m Model) buildPreviewView() string {
	return m.renderer.Render(m.tree)
}

func (m Model) buildHelpLine() string {
	keys := []string{
		"↑↓:move",
		"Tab:indent",
		"S-Tab:unindent",
		"Enter:new",
		"C-d:delete",
		"C-c:copy",
		"C-q:quit",
	}
	return helpStyle.Render(fmt.Sprintf(" %s ", strings.Join(keys, " │ ")))
}
