package ui_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"gocheck-net/internal/checker"
	"gocheck-net/internal/ui"
)

func TestRenderSpeedBar(t *testing.T) {
	theme := ui.Themes[0] // Default theme

	t.Run("Checking state shows probing animation placeholder", func(t *testing.T) {
		res := ui.RenderSpeedBar(0, true, false, theme)
		if !strings.Contains(res, "probing") {
			t.Errorf("Expected 'probing' in speed bar, got %q", res)
		}
	})

	t.Run("Failed probe shows failed indicator", func(t *testing.T) {
		res := ui.RenderSpeedBar(0, false, false, theme)
		if !strings.Contains(res, "failed") {
			t.Errorf("Expected 'failed' in speed bar, got %q", res)
		}
	})

	t.Run("Fast latency renders blazing emoji", func(t *testing.T) {
		res := ui.RenderSpeedBar(45*time.Millisecond, false, true, theme)
		if !strings.Contains(res, "⚡") {
			t.Errorf("Expected '⚡' in fast speed bar, got %q", res)
		}
	})

	t.Run("Medium latency renders walking emoji", func(t *testing.T) {
		res := ui.RenderSpeedBar(450*time.Millisecond, false, true, theme)
		if !strings.Contains(res, "🚶") {
			t.Errorf("Expected '🚶' in medium speed bar, got %q", res)
		}
	})

	t.Run("Slow latency renders snail emoji", func(t *testing.T) {
		res := ui.RenderSpeedBar(950*time.Millisecond, false, true, theme)
		if !strings.Contains(res, "🐌") {
			t.Errorf("Expected '🐌' in slow speed bar, got %q", res)
		}
	})
}

func TestRenderHealthBar(t *testing.T) {
	theme := ui.Themes[0]

	t.Run("100% health renders full percentage", func(t *testing.T) {
		res := ui.RenderHealthBar(100, theme)
		if !strings.Contains(res, "100%") {
			t.Errorf("Expected '100%%' in health bar, got %q", res)
		}
	})

	t.Run("50% health renders 50%", func(t *testing.T) {
		res := ui.RenderHealthBar(50, theme)
		if !strings.Contains(res, "50%") {
			t.Errorf("Expected '50%%' in health bar, got %q", res)
		}
	})

	t.Run("0% health renders 0%", func(t *testing.T) {
		res := ui.RenderHealthBar(0, theme)
		if !strings.Contains(res, "0%") {
			t.Errorf("Expected '0%%' in health bar, got %q", res)
		}
	})
}

func TestModelKeyboardInteractions(t *testing.T) {
	targets := []checker.Target{
		checker.ParseTarget("https://google.com"),
		checker.ParseTarget("https://github.com"),
		checker.ParseTarget("1.1.1.1:53"),
	}

	model := ui.New(targets, 2*time.Second)

	// Test cursor movement down
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m := newModel.(ui.Model)
	if m.Cursor != 1 {
		t.Errorf("Expected cursor at 1 after 'j', got %d", m.Cursor)
	}

	// Test cursor movement up
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = newModel.(ui.Model)
	if m.Cursor != 0 {
		t.Errorf("Expected cursor at 0 after 'k', got %d", m.Cursor)
	}

	// Test theme cycling
	initialTheme := m.CurrentTheme
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = newModel.(ui.Model)
	if m.CurrentTheme != (initialTheme+1)%len(ui.Themes) {
		t.Errorf("Expected theme %d, got %d", (initialTheme+1)%len(ui.Themes), m.CurrentTheme)
	}

	// Test audio alert toggle
	initialAlert := m.AudioAlert
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m = newModel.(ui.Model)
	if m.AudioAlert == initialAlert {
		t.Errorf("Expected AudioAlert to toggle from %v, got %v", initialAlert, m.AudioAlert)
	}

	// Test deletion
	initialLen := len(m.Targets)
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	m = newModel.(ui.Model)
	if len(m.Targets) != initialLen-1 {
		t.Errorf("Expected %d targets after deletion, got %d", initialLen-1, len(m.Targets))
	}

	// Test entering input mode
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	m = newModel.(ui.Model)
	if !m.InputMode {
		t.Errorf("Expected InputMode to be true after pressing '+'")
	}

	// Test exiting input mode with escape
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newModel.(ui.Model)
	if m.InputMode {
		t.Errorf("Expected InputMode to be false after pressing Esc")
	}
}

func TestViewRendering(t *testing.T) {
	targets := []checker.Target{
		checker.ParseTarget("https://google.com"),
	}

	model := ui.New(targets, 2*time.Second)
	viewOutput := model.View()

	if !strings.Contains(viewOutput, "GOCHECK-NET") {
		t.Errorf("View missing title header")
	}
	if !strings.Contains(viewOutput, "SPEED GAUGE") {
		t.Errorf("View missing SPEED GAUGE column")
	}
	if !strings.Contains(viewOutput, "Theme (") {
		t.Errorf("View missing Theme hotkey guide")
	}
}
