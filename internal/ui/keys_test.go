package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()

	if len(km.Up) == 0 {
		t.Error("expected Up keys to be defined")
	}
	if len(km.Down) == 0 {
		t.Error("expected Down keys to be defined")
	}
	if len(km.Indent) == 0 {
		t.Error("expected Indent keys to be defined")
	}
	if len(km.Unindent) == 0 {
		t.Error("expected Unindent keys to be defined")
	}
	if len(km.Enter) == 0 {
		t.Error("expected Enter keys to be defined")
	}
	if len(km.Delete) == 0 {
		t.Error("expected Delete keys to be defined")
	}
	if len(km.Copy) == 0 {
		t.Error("expected Copy keys to be defined")
	}
	if len(km.Quit) == 0 {
		t.Error("expected Quit keys to be defined")
	}
}

func TestMatchesUp(t *testing.T) {
	km := DefaultKeyMap()

	tests := []struct {
		key      tea.KeyMsg
		expected bool
	}{
		{tea.KeyMsg{Type: tea.KeyUp}, true},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}, true},
		{tea.KeyMsg{Type: tea.KeyDown}, false},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}, false},
	}

	for _, tt := range tests {
		result := matches(tt.key, km.Up)
		if result != tt.expected {
			t.Errorf("matches(%v, Up) = %v, expected %v", tt.key, result, tt.expected)
		}
	}
}

func TestMatchesDown(t *testing.T) {
	km := DefaultKeyMap()

	tests := []struct {
		key      tea.KeyMsg
		expected bool
	}{
		{tea.KeyMsg{Type: tea.KeyDown}, true},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}, true},
		{tea.KeyMsg{Type: tea.KeyUp}, false},
	}

	for _, tt := range tests {
		result := matches(tt.key, km.Down)
		if result != tt.expected {
			t.Errorf("matches(%v, Down) = %v, expected %v", tt.key, result, tt.expected)
		}
	}
}

func TestMatchesIndent(t *testing.T) {
	km := DefaultKeyMap()

	msg := tea.KeyMsg{Type: tea.KeyTab}
	if !matches(msg, km.Indent) {
		t.Error("Tab should match Indent")
	}
}

func TestMatchesUnindent(t *testing.T) {
	km := DefaultKeyMap()

	msg := tea.KeyMsg{Type: tea.KeyShiftTab}
	if !matches(msg, km.Unindent) {
		t.Error("Shift+Tab should match Unindent")
	}
}

func TestMatchesEnter(t *testing.T) {
	km := DefaultKeyMap()

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	if !matches(msg, km.Enter) {
		t.Error("Enter should match Enter keys")
	}
}

func TestMatchesDelete(t *testing.T) {
	km := DefaultKeyMap()

	msg := tea.KeyMsg{Type: tea.KeyCtrlD}
	if !matches(msg, km.Delete) {
		t.Error("Ctrl+D should match Delete")
	}
}

func TestMatchesCopy(t *testing.T) {
	km := DefaultKeyMap()

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	if !matches(msg, km.Copy) {
		t.Error("Ctrl+C should match Copy")
	}
}

func TestMatchesQuit(t *testing.T) {
	km := DefaultKeyMap()

	tests := []struct {
		key      tea.KeyMsg
		expected bool
	}{
		{tea.KeyMsg{Type: tea.KeyCtrlQ}, true},
		{tea.KeyMsg{Type: tea.KeyEscape}, true},
		{tea.KeyMsg{Type: tea.KeyEnter}, false},
	}

	for _, tt := range tests {
		result := matches(tt.key, km.Quit)
		if result != tt.expected {
			t.Errorf("matches(%v, Quit) = %v, expected %v", tt.key, result, tt.expected)
		}
	}
}

func TestMatchesEmpty(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	if matches(msg, []string{}) {
		t.Error("should not match empty key list")
	}
}
