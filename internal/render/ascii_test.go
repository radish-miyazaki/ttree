package render

import (
	"strings"
	"testing"

	"github.com/radish-miyazaki/ttree/internal/tree"
)

func TestDefaultStyle(t *testing.T) {
	style := DefaultStyle()

	if style.Branch != "├── " {
		t.Errorf("unexpected Branch: %q", style.Branch)
	}
	if style.LastBranch != "└── " {
		t.Errorf("unexpected LastBranch: %q", style.LastBranch)
	}
	if style.Vertical != "│   " {
		t.Errorf("unexpected Vertical: %q", style.Vertical)
	}
	if style.Space != "    " {
		t.Errorf("unexpected Space: %q", style.Space)
	}
}

func TestRenderSingleNode(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil
	tr.Root.AddChild(tree.NewNode("item"))

	r := NewRenderer()
	output := r.Render(tr)

	expected := "└── item\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRenderMultipleSiblings(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil
	tr.Root.AddChild(tree.NewNode("item1"))
	tr.Root.AddChild(tree.NewNode("item2"))
	tr.Root.AddChild(tree.NewNode("item3"))

	r := NewRenderer()
	output := r.Render(tr)

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "├── item1" {
		t.Errorf("line 0: expected '├── item1', got %q", lines[0])
	}
	if lines[1] != "├── item2" {
		t.Errorf("line 1: expected '├── item2', got %q", lines[1])
	}
	if lines[2] != "└── item3" {
		t.Errorf("line 2: expected '└── item3', got %q", lines[2])
	}
}

func TestRenderNestedNodes(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil

	parent := tree.NewNode("parent")
	child := tree.NewNode("child")
	parent.AddChild(child)
	tr.Root.AddChild(parent)

	r := NewRenderer()
	output := r.Render(tr)

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "└── parent" {
		t.Errorf("line 0: expected '└── parent', got %q", lines[0])
	}
	if lines[1] != "    └── child" {
		t.Errorf("line 1: expected '    └── child', got %q", lines[1])
	}
}

func TestRenderComplexTree(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil

	// Build tree:
	// ├── folder1
	// │   ├── file1
	// │   └── file2
	// └── folder2
	//     └── file3

	folder1 := tree.NewNode("folder1")
	folder1.AddChild(tree.NewNode("file1"))
	folder1.AddChild(tree.NewNode("file2"))

	folder2 := tree.NewNode("folder2")
	folder2.AddChild(tree.NewNode("file3"))

	tr.Root.AddChild(folder1)
	tr.Root.AddChild(folder2)

	r := NewRenderer()
	output := r.Render(tr)

	expected := `├── folder1
│   ├── file1
│   └── file2
└── folder2
    └── file3
`
	if output != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, output)
	}
}

func TestRenderEmptyText(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil
	tr.Root.AddChild(tree.NewNode(""))

	r := NewRenderer()
	output := r.Render(tr)

	// Empty text should render as space
	expected := "└──  \n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRenderCollapsedNode(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil

	parent := tree.NewNode("parent")
	parent.AddChild(tree.NewNode("child"))
	parent.Expanded = false
	tr.Root.AddChild(parent)

	r := NewRenderer()
	output := r.Render(tr)

	// Child should not be rendered
	expected := "└── parent\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestRenderLines(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil
	tr.Root.AddChild(tree.NewNode("item1"))
	tr.Root.AddChild(tree.NewNode("item2"))

	r := NewRenderer()
	lines := r.RenderLines(tr)

	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestRenderLinesEmpty(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil

	r := NewRenderer()
	lines := r.RenderLines(tr)

	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(lines))
	}
}

func TestRenderDeepNesting(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil

	// Create deep nesting: root -> a -> b -> c -> d
	a := tree.NewNode("a")
	b := tree.NewNode("b")
	c := tree.NewNode("c")
	d := tree.NewNode("d")

	a.AddChild(b)
	b.AddChild(c)
	c.AddChild(d)
	tr.Root.AddChild(a)

	r := NewRenderer()
	output := r.Render(tr)

	expected := `└── a
    └── b
        └── c
            └── d
`
	if output != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, output)
	}
}

func TestRenderMixedStructure(t *testing.T) {
	tr := tree.NewTree()
	tr.Root.Children = nil

	// Build:
	// ├── a
	// │   └── a1
	// ├── b
	// └── c
	//     ├── c1
	//     └── c2

	a := tree.NewNode("a")
	a.AddChild(tree.NewNode("a1"))

	b := tree.NewNode("b")

	c := tree.NewNode("c")
	c.AddChild(tree.NewNode("c1"))
	c.AddChild(tree.NewNode("c2"))

	tr.Root.AddChild(a)
	tr.Root.AddChild(b)
	tr.Root.AddChild(c)

	r := NewRenderer()
	output := r.Render(tr)

	expected := `├── a
│   └── a1
├── b
└── c
    ├── c1
    └── c2
`
	if output != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, output)
	}
}
