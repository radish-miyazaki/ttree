package ui

import "github.com/charmbracelet/bubbletea"

// KeyMap defines all key bindings
type KeyMap struct {
	Up        []string
	Down      []string
	Left      []string
	Right     []string
	Indent    []string
	Unindent  []string
	Enter     []string
	Delete    []string
	Copy      []string
	Quit      []string
	Help      []string
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:       []string{"up", "k"},
		Down:     []string{"down", "j"},
		Left:     []string{"left"},
		Right:    []string{"right"},
		Indent:   []string{"tab"},
		Unindent: []string{"shift+tab"},
		Enter:    []string{"enter"},
		Delete:   []string{"ctrl+d", "ctrl+backspace"},
		Copy:     []string{"ctrl+c"},
		Quit:     []string{"ctrl+q", "esc"},
		Help:     []string{"ctrl+?", "f1"},
	}
}

// matches checks if a key message matches any of the given keys
func matches(msg tea.KeyMsg, keys []string) bool {
	for _, k := range keys {
		if msg.String() == k {
			return true
		}
	}
	return false
}
