package ui

import (
	"time"

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
	Targets  []TargetItem
	Timeout  time.Duration
	Quitting bool
}

// checkDoneMsg is delivered to Update when an async probe completes.
type checkDoneMsg struct {
	Index  int
	Result checker.Result
}

// New creates and initializes a Model with target items.
func New(targets []checker.Target, timeout time.Duration) Model {
	items := make([]TargetItem, len(targets))
	for i, t := range targets {
		items[i] = TargetItem{
			Target:   t,
			Checking: true,
			Result:   nil,
		}
	}

	return Model{
		Targets: items,
		Timeout: timeout,
	}
}

// Compile-time check to ensure Model implements tea.Model interface.
var _ tea.Model = Model{}
