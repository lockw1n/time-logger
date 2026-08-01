package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// textInput is a deliberately tiny single-line editor: a value plus a cursor
// position. The plan rules out the bubbles component kit unless genuinely needed,
// and the two forms here (a duration field, an add-entry form) only ever need
// insert, backspace, delete and horizontal movement — a handful of lines, not a
// dependency.
type textInput struct {
	value string
	pos   int // cursor index in runes [0..len(runes)]
}

func newTextInput(value string) textInput {
	return textInput{value: value, pos: len([]rune(value))}
}

// trimmedEmpty reports whether the field holds only whitespace — the signal the
// duration editor uses to mean "delete this entry" rather than set a new value.
func (t textInput) trimmedEmpty() bool {
	return strings.TrimSpace(t.value) == ""
}

// update applies a key to the input and reports whether it consumed the key, so
// the caller can let unhandled keys (tab, enter, esc) fall through to the form.
func (t *textInput) update(msg tea.KeyMsg) (handled bool) {
	runes := []rune(t.value)
	switch msg.Type {
	case tea.KeyRunes, tea.KeySpace:
		ins := msg.Runes
		if msg.Type == tea.KeySpace {
			ins = []rune{' '}
		}
		runes = append(runes[:t.pos], append(append([]rune{}, ins...), runes[t.pos:]...)...)
		t.pos += len(ins)
	case tea.KeyBackspace:
		if t.pos > 0 {
			runes = append(runes[:t.pos-1], runes[t.pos:]...)
			t.pos--
		}
	case tea.KeyDelete:
		if t.pos < len(runes) {
			runes = append(runes[:t.pos], runes[t.pos+1:]...)
		}
	case tea.KeyLeft:
		if t.pos > 0 {
			t.pos--
		}
	case tea.KeyRight:
		if t.pos < len(runes) {
			t.pos++
		}
	case tea.KeyHome:
		t.pos = 0
	case tea.KeyEnd:
		t.pos = len(runes)
	default:
		return false
	}
	t.value = string(runes)
	return true
}

// render draws the value with a block cursor when focused. It appends a cursor
// cell past the end of the text so an empty focused field is still visible.
func (t textInput) render(focused bool) string {
	if !focused {
		if t.value == "" {
			return " "
		}
		return t.value
	}
	runes := []rune(t.value)
	if t.pos >= len(runes) {
		return string(runes) + cursorStyle.Render(" ")
	}
	return string(runes[:t.pos]) + cursorStyle.Render(string(runes[t.pos])) + string(runes[t.pos+1:])
}
