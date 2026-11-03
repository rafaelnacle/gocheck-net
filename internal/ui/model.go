package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"gocheck-net/internal/checker"
)

// TargetItem encapsulates a single target and its current status in the UI.
type TargetItem struct {
	Target   checker.Target
	Checking bool
	Result   *checker.Result
}

// Model represents the Bubble Tea application state.
type Model struct {
	Targets       []TargetItem
	Timeout       time.Duration
	Cursor        int             // Selected row in table
	CurrentTheme  int             // Index in Themes slice
	AudioAlert    bool            // Play terminal bell on downtime
	InputMode     bool            // True when typing new target
	InspectorMode bool            // True when viewing detailed inspection modal
	TextInput     textinput.Model // Bubbles textinput widget
	FlashNotice   string          // Brief notification message
	Quitting      bool
}

// checkDoneMsg is delivered to Update when an async probe completes.
type checkDoneMsg struct {
	Index  int
	Result checker.Result
}

// New creates and initializes a Model with target items, bubbles textinput, and default settings.
func New(targets []checker.Target, timeout time.Duration) Model {
	items := make([]TargetItem, len(targets))
	for i, t := range targets {
		items[i] = TargetItem{
			Target:   t,
			Checking: true,
			Result:   nil,
		}
	}

	ti := textinput.New()
	ti.Placeholder = "e.g. reddit.com, https://api.github.com, or 1.1.1.1:53"
	ti.Prompt = " ❯ Add Target: "
	ti.CharLimit = 120

	return Model{
		Targets:       items,
		Timeout:       timeout,
		Cursor:        0,
		CurrentTheme:  0,
		AudioAlert:    true, // Enabled by default for outages
		InputMode:     false,
		InspectorMode: false,
		TextInput:     ti,
		Quitting:      false,
	}
}

// Compile-time check to ensure Model implements tea.Model interface.
var _ tea.Model = Model{}
