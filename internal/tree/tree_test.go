package tree

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	node := NewNode("test")

	if node.Text != "test" {
		t.Errorf("expected text 'test', got '%s'", node.Text)
	}
	if node.ID == "" {
		t.Error("expected non-empty ID")
	}
	if len(node.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(node.Children))
	}
	if !node.Expanded {
		t.Error("expected node to be expanded by default")
	}
}

func TestAddChild(t *testing.T) {
	parent := NewNode("parent")
	child := NewNode("child")

	parent.AddChild(child)

	if len(parent.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(parent.Children))
	}
	if parent.Children[0] != child {
		t.Error("child not added correctly")
	}
	if child.Parent != parent {
		t.Error("parent reference not set")
	}
}

func TestAddChildAt(t *testing.T) {
	parent := NewNode("parent")
	child1 := NewNode("child1")
	child2 := NewNode("child2")
	child3 := NewNode("child3")

	parent.AddChild(child1)
	parent.AddChild(child3)
	parent.AddChildAt(child2, 1)

	if len(parent.Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(parent.Children))
	}
	if parent.Children[0].Text != "child1" {
		t.Errorf("expected child1 at index 0, got %s", parent.Children[0].Text)
	}
	if parent.Children[1].Text != "child2" {
		t.Errorf("expected child2 at index 1, got %s", parent.Children[1].Text)
	}
	if parent.Children[2].Text != "child3" {
		t.Errorf("expected child3 at index 2, got %s", parent.Children[2].Text)
	}
}

func TestRemoveChild(t *testing.T) {
	parent := NewNode("parent")
	child1 := NewNode("child1")
	child2 := NewNode("child2")

	parent.AddChild(child1)
	parent.AddChild(child2)

	removed := parent.RemoveChild(child1)

	if !removed {
		t.Error("expected RemoveChild to return true")
	}
	if len(parent.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(parent.Children))
	}
	if parent.Children[0] != child2 {
		t.Error("wrong child remaining")
	}
	if child1.Parent != nil {
		t.Error("removed child should have nil parent")
	}
}

func TestRemoveChildNotFound(t *testing.T) {
	parent := NewNode("parent")
	child := NewNode("child")
	other := NewNode("other")

	parent.AddChild(child)

	removed := parent.RemoveChild(other)

	if removed {
		t.Error("expected RemoveChild to return false for non-existent child")
	}
}

func TestDepth(t *testing.T) {
	root := NewNode("root")
	child := NewNode("child")
	grandchild := NewNode("grandchild")

	root.AddChild(child)
	child.AddChild(grandchild)

	if root.Depth() != 0 {
		t.Errorf("expected root depth 0, got %d", root.Depth())
	}
	if child.Depth() != 1 {
		t.Errorf("expected child depth 1, got %d", child.Depth())
	}
	if grandchild.Depth() != 2 {
		t.Errorf("expected grandchild depth 2, got %d", grandchild.Depth())
	}
}

func TestIndex(t *testing.T) {
	parent := NewNode("parent")
	child1 := NewNode("child1")
	child2 := NewNode("child2")
	child3 := NewNode("child3")

	parent.AddChild(child1)
	parent.AddChild(child2)
	parent.AddChild(child3)

	if child1.Index() != 0 {
		t.Errorf("expected index 0, got %d", child1.Index())
	}
	if child2.Index() != 1 {
		t.Errorf("expected index 1, got %d", child2.Index())
	}
	if child3.Index() != 2 {
		t.Errorf("expected index 2, got %d", child3.Index())
	}
}

func TestIsLastChild(t *testing.T) {
	parent := NewNode("parent")
	child1 := NewNode("child1")
	child2 := NewNode("child2")

	parent.AddChild(child1)
	parent.AddChild(child2)

	if child1.IsLastChild() {
		t.Error("child1 should not be last child")
	}
	if !child2.IsLastChild() {
		t.Error("child2 should be last child")
	}
}

func TestNewTree(t *testing.T) {
	tree := NewTree()

	if tree.Root == nil {
		t.Error("expected non-nil root")
	}
	if len(tree.Root.Children) != 1 {
		t.Errorf("expected 1 initial child, got %d", len(tree.Root.Children))
	}
}

func TestFlattenVisible(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil // Clear default child

	child1 := NewNode("child1")
	child2 := NewNode("child2")
	grandchild := NewNode("grandchild")

	tree.Root.AddChild(child1)
	tree.Root.AddChild(child2)
	child1.AddChild(grandchild)

	nodes := tree.FlattenVisible()

	if len(nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(nodes))
	}
	if nodes[0].Text != "child1" {
		t.Errorf("expected first node 'child1', got '%s'", nodes[0].Text)
	}
	if nodes[1].Text != "grandchild" {
		t.Errorf("expected second node 'grandchild', got '%s'", nodes[1].Text)
	}
	if nodes[2].Text != "child2" {
		t.Errorf("expected third node 'child2', got '%s'", nodes[2].Text)
	}
}

func TestFlattenVisibleCollapsed(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil

	child1 := NewNode("child1")
	grandchild := NewNode("grandchild")

	tree.Root.AddChild(child1)
	child1.AddChild(grandchild)
	child1.Expanded = false

	nodes := tree.FlattenVisible()

	if len(nodes) != 1 {
		t.Errorf("expected 1 node when collapsed, got %d", len(nodes))
	}
}

func TestIndent(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil

	child1 := NewNode("child1")
	child2 := NewNode("child2")

	tree.Root.AddChild(child1)
	tree.Root.AddChild(child2)

	// child2 should become child of child1
	result := tree.Indent(child2)

	if !result {
		t.Error("expected Indent to return true")
	}
	if len(tree.Root.Children) != 1 {
		t.Errorf("expected 1 root child, got %d", len(tree.Root.Children))
	}
	if len(child1.Children) != 1 {
		t.Errorf("expected child1 to have 1 child, got %d", len(child1.Children))
	}
	if child1.Children[0] != child2 {
		t.Error("child2 should be child of child1")
	}
}

func TestIndentFirstChild(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil

	child1 := NewNode("child1")
	tree.Root.AddChild(child1)

	// Cannot indent first child (no previous sibling)
	result := tree.Indent(child1)

	if result {
		t.Error("expected Indent to return false for first child")
	}
}

func TestUnindent(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil

	child1 := NewNode("child1")
	grandchild := NewNode("grandchild")

	tree.Root.AddChild(child1)
	child1.AddChild(grandchild)

	// grandchild should become sibling of child1
	result := tree.Unindent(grandchild)

	if !result {
		t.Error("expected Unindent to return true")
	}
	if len(tree.Root.Children) != 2 {
		t.Errorf("expected 2 root children, got %d", len(tree.Root.Children))
	}
	if tree.Root.Children[1] != grandchild {
		t.Error("grandchild should be second root child")
	}
}

func TestUnindentRootChild(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil

	child1 := NewNode("child1")
	tree.Root.AddChild(child1)

	// Cannot unindent root's direct child
	result := tree.Unindent(child1)

	if result {
		t.Error("expected Unindent to return false for root's child")
	}
}

func TestInsertAfter(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil

	child1 := NewNode("child1")
	child2 := NewNode("child2")
	newNode := NewNode("new")

	tree.Root.AddChild(child1)
	tree.Root.AddChild(child2)

	tree.InsertAfter(child1, newNode)

	if len(tree.Root.Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(tree.Root.Children))
	}
	if tree.Root.Children[1] != newNode {
		t.Error("new node should be at index 1")
	}
}

func TestDelete(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil

	child1 := NewNode("child1")
	child2 := NewNode("child2")

	tree.Root.AddChild(child1)
	tree.Root.AddChild(child2)

	nextFocus := tree.Delete(child1)

	if len(tree.Root.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(tree.Root.Children))
	}
	if nextFocus != child2 {
		t.Error("expected next focus to be child2")
	}
}

func TestDeleteOnlyChild(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil

	child := NewNode("child")
	child.Text = "some text"
	tree.Root.AddChild(child)

	nextFocus := tree.Delete(child)

	// Should not delete, but clear text
	if len(tree.Root.Children) != 1 {
		t.Errorf("expected 1 child (not deleted), got %d", len(tree.Root.Children))
	}
	if child.Text != "" {
		t.Errorf("expected empty text, got '%s'", child.Text)
	}
	if nextFocus != child {
		t.Error("expected next focus to be the same child")
	}
}

func TestDeleteReturnsParent(t *testing.T) {
	tree := NewTree()
	tree.Root.Children = nil

	parent := NewNode("parent")
	child := NewNode("child")

	tree.Root.AddChild(parent)
	parent.AddChild(child)

	nextFocus := tree.Delete(child)

	if nextFocus != parent {
		t.Error("expected next focus to be parent when last child deleted")
	}
}
