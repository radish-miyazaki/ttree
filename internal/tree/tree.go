package tree

import "github.com/google/uuid"

// Node represents a single node in the tree
type Node struct {
	ID       string
	Text     string
	Children []*Node
	Parent   *Node
	Expanded bool
}

// NewNode creates a new node with a unique ID
func NewNode(text string) *Node {
	return &Node{
		ID:       uuid.New().String()[:8],
		Text:     text,
		Children: make([]*Node, 0),
		Expanded: true,
	}
}

// AddChild adds a child node
func (n *Node) AddChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

// AddChildAt adds a child node at a specific index
func (n *Node) AddChildAt(child *Node, index int) {
	child.Parent = n
	if index >= len(n.Children) {
		n.Children = append(n.Children, child)
		return
	}
	n.Children = append(n.Children[:index+1], n.Children[index:]...)
	n.Children[index] = child
}

// RemoveChild removes a child node
func (n *Node) RemoveChild(child *Node) bool {
	for i, c := range n.Children {
		if c.ID == child.ID {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			child.Parent = nil
			return true
		}
	}
	return false
}

// Depth returns the depth of this node (root = 0)
func (n *Node) Depth() int {
	depth := 0
	current := n.Parent
	for current != nil {
		depth++
		current = current.Parent
	}
	return depth
}

// Index returns the index of this node among its siblings
func (n *Node) Index() int {
	if n.Parent == nil {
		return 0
	}
	for i, c := range n.Parent.Children {
		if c.ID == n.ID {
			return i
		}
	}
	return -1
}

// IsLastChild returns true if this node is the last child of its parent
func (n *Node) IsLastChild() bool {
	if n.Parent == nil {
		return true
	}
	return n.Index() == len(n.Parent.Children)-1
}

// Tree represents the entire tree structure
type Tree struct {
	Root *Node
}

// NewTree creates a new tree with a root node
func NewTree() *Tree {
	root := NewNode("root")
	root.AddChild(NewNode(""))
	return &Tree{Root: root}
}

// FlattenVisible returns all visible nodes in order (for display)
func (t *Tree) FlattenVisible() []*Node {
	var result []*Node
	t.flattenNode(t.Root, &result, true)
	return result
}

func (t *Tree) flattenNode(n *Node, result *[]*Node, skipRoot bool) {
	if !skipRoot {
		*result = append(*result, n)
	}
	if n.Expanded || skipRoot {
		for _, child := range n.Children {
			t.flattenNode(child, result, false)
		}
	}
}

// Indent moves a node to be a child of its previous sibling
func (t *Tree) Indent(n *Node) bool {
	if n.Parent == nil {
		return false
	}
	idx := n.Index()
	if idx <= 0 {
		return false
	}
	prevSibling := n.Parent.Children[idx-1]
	n.Parent.RemoveChild(n)
	prevSibling.AddChild(n)
	prevSibling.Expanded = true
	return true
}

// Unindent moves a node to be a sibling of its parent
func (t *Tree) Unindent(n *Node) bool {
	if n.Parent == nil || n.Parent.Parent == nil {
		return false
	}
	grandparent := n.Parent.Parent
	parentIdx := n.Parent.Index()
	n.Parent.RemoveChild(n)
	grandparent.AddChildAt(n, parentIdx+1)
	return true
}

// InsertAfter inserts a new node after the given node
func (t *Tree) InsertAfter(n *Node, newNode *Node) {
	if n.Parent == nil {
		return
	}
	idx := n.Index()
	n.Parent.AddChildAt(newNode, idx+1)
}

// InsertChild inserts a new node as the first child
func (t *Tree) InsertChild(n *Node, newNode *Node) {
	n.AddChildAt(newNode, 0)
	n.Expanded = true
}

// Delete removes a node and returns the node to focus next
func (t *Tree) Delete(n *Node) *Node {
	if n.Parent == nil {
		return nil
	}
	parent := n.Parent
	idx := n.Index()

	// Don't delete if it's the only node
	if parent == t.Root && len(parent.Children) == 1 {
		n.Text = ""
		return n
	}

	parent.RemoveChild(n)

	// Return next focus target
	if idx > 0 {
		return parent.Children[idx-1]
	}
	if len(parent.Children) > 0 {
		return parent.Children[0]
	}
	if parent != t.Root {
		return parent
	}
	return nil
}
