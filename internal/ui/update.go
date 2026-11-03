package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"gocheck-net/internal/checker"
)

// Init executes immediately when the program starts.
// It returns a batch command to probe all targets concurrently.
func (m Model) Init() tea.Cmd {
	return m.checkAllCmd()
}

// Update processes messages (user input, completed probes) and returns the updated model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit

		case "r":
			// Mark all targets as currently checking and trigger probes
			for i := range m.Targets {
				m.Targets[i].Checking = true
			}
			return m, m.checkAllCmd()
		}

	case checkDoneMsg:
		if msg.Index >= 0 && msg.Index < len(m.Targets) {
			m.Targets[msg.Index].Result = &msg.Result
			m.Targets[msg.Index].Checking = false
		}
		return m, nil
	}

	return m, nil
}

// checkAllCmd generates a batched command that probes all targets in parallel.
func (m Model) checkAllCmd() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Targets))
	for i, item := range m.Targets {
		cmds[i] = checkTargetCmd(i, item.Target, m.Timeout)
	}
	return tea.Batch(cmds...)
}

// checkTargetCmd produces a tea.Cmd (an async worker function) that runs checker.Check.
func checkTargetCmd(index int, target checker.Target, timeout time.Duration) tea.Cmd {
	return func() tea.Msg {
		// Run probe with background context and configured timeout
		res := checker.Check(context.Background(), target, timeout)
		return checkDoneMsg{
			Index:  index,
			Result: res,
		}
	}
}
