package ui

import "github.com/radish-miyazaki/ttree/internal/tree"

// maxHistory caps the number of undo snapshots kept in memory.
const maxHistory = 100

// snapshot captures the tree and cursor position at a point in time.
type snapshot struct {
	tree   *tree.Tree
	cursor int
}

// history holds the undo and redo stacks of tree snapshots.
type history struct {
	undo []snapshot
	redo []snapshot
}

// record pushes a snapshot onto the undo stack and clears the redo stack.
func (h *history) record(s snapshot) {
	h.undo = append(h.undo, s)
	if len(h.undo) > maxHistory {
		h.undo = h.undo[len(h.undo)-maxHistory:]
	}
	h.redo = nil
}

// undoTo pops the latest undo snapshot, pushing current onto the redo stack.
func (h *history) undoTo(current snapshot) (snapshot, bool) {
	if len(h.undo) == 0 {
		return snapshot{}, false
	}
	s := h.undo[len(h.undo)-1]
	h.undo = h.undo[:len(h.undo)-1]
	h.redo = append(h.redo, current)
	return s, true
}

// redoTo pops the latest redo snapshot, pushing current onto the undo stack.
func (h *history) redoTo(current snapshot) (snapshot, bool) {
	if len(h.redo) == 0 {
		return snapshot{}, false
	}
	s := h.redo[len(h.redo)-1]
	h.redo = h.redo[:len(h.redo)-1]
	h.undo = append(h.undo, current)
	return s, true
}
