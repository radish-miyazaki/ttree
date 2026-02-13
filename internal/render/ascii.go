package render

import (
	"strings"

	"github.com/radish-miyazaki/ttree/internal/tree"
)

// Style defines the characters used for tree rendering
type Style struct {
	Branch     string // ├──
	LastBranch string // └──
	Vertical   string // │
	Space      string //
}

// DefaultStyle returns the default ASCII tree style
func DefaultStyle() Style {
	return Style{
		Branch:     "├── ",
		LastBranch: "└── ",
		Vertical:   "│   ",
		Space:      "    ",
	}
}

// Renderer renders tree structures to ASCII art
type Renderer struct {
	Style Style
}

// NewRenderer creates a new ASCII renderer
func NewRenderer() *Renderer {
	return &Renderer{Style: DefaultStyle()}
}

// Render renders the entire tree to a string
func (r *Renderer) Render(t *tree.Tree) string {
	var sb strings.Builder
	for i, child := range t.Root.Children {
		isLast := i == len(t.Root.Children)-1
		r.renderNode(&sb, child, "", isLast)
	}
	return sb.String()
}

// RenderLines renders the tree and returns individual lines
func (r *Renderer) RenderLines(t *tree.Tree) []string {
	output := r.Render(t)
	if output == "" {
		return []string{}
	}
	lines := strings.Split(output, "\n")
	// Remove trailing empty line
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func (r *Renderer) renderNode(sb *strings.Builder, n *tree.Node, prefix string, isLast bool) {
	// Choose branch character
	branch := r.Style.Branch
	if isLast {
		branch = r.Style.LastBranch
	}

	// Write current node
	text := n.Text
	if text == "" {
		text = " "
	}
	sb.WriteString(prefix + branch + text + "\n")

	// Calculate prefix for children
	childPrefix := prefix
	if isLast {
		childPrefix += r.Style.Space
	} else {
		childPrefix += r.Style.Vertical
	}

	// Render children if expanded
	if n.Expanded {
		for i, child := range n.Children {
			childIsLast := i == len(n.Children)-1
			r.renderNode(sb, child, childPrefix, childIsLast)
		}
	}
}

// GetPrefix returns the prefix string for a node at given depth
func (r *Renderer) GetPrefix(ancestors []bool) string {
	var prefix strings.Builder
	for i, isLast := range ancestors[:len(ancestors)-1] {
		if i == 0 {
			continue
		}
		if isLast {
			prefix.WriteString(r.Style.Space)
		} else {
			prefix.WriteString(r.Style.Vertical)
		}
	}
	return prefix.String()
}
