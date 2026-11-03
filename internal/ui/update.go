package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"gocheck-net/internal/checker"
)

// Init executes immediately when the program starts.
// It triggers probes for all targets and starts text input blinking.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.checkAllCmd(), textinput.Blink)
}

// ringBellCmd triggers the terminal audio bell character.
func ringBellCmd() tea.Msg {
	fmt.Print("\a")
	return nil
}

// Update processes messages (user input, completed probes) and returns the updated model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If user is currently typing a new target in the interactive input box
	if m.InputMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				val := strings.TrimSpace(m.TextInput.Value())
				if val != "" {
					newTarget := checker.ParseTarget(val)
					m.Targets = append(m.Targets, TargetItem{
						Target:   newTarget,
						Checking: true,
						Result:   nil,
					})
					newIndex := len(m.Targets) - 1
					m.Cursor = newIndex
					m.TextInput.Reset()
					m.InputMode = false
					m.FlashNotice = fmt.Sprintf("Added %s", newTarget.Raw)
					return m, checkTargetCmd(newIndex, newTarget, m.Timeout)
				}
				m.InputMode = false
				m.TextInput.Reset()
				return m, nil

			case "esc":
				m.InputMode = false
				m.TextInput.Reset()
				return m, nil
			}

		case checkDoneMsg:
			return m.handleCheckDone(msg)
		}

		var cmd tea.Cmd
		m.TextInput, cmd = m.TextInput.Update(msg)
		return m, cmd
	}

	// Normal navigation and command mode
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit

		case "esc":
			if m.FlashNotice != "" {
				m.FlashNotice = ""
				return m, nil
			}
			m.Quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil

		case "down", "j":
			if m.Cursor < len(m.Targets)-1 {
				m.Cursor++
			}
			return m, nil

		case "+", "a":
			m.InputMode = true
			m.TextInput.Reset()
			cmd := m.TextInput.Focus()
			return m, cmd

		case "x", "d":
			if len(m.Targets) > 0 && m.Cursor < len(m.Targets) {
				removed := m.Targets[m.Cursor].Target.Raw
				m.Targets = append(m.Targets[:m.Cursor], m.Targets[m.Cursor+1:]...)
				if m.Cursor >= len(m.Targets) && m.Cursor > 0 {
					m.Cursor = len(m.Targets) - 1
				}
				m.FlashNotice = fmt.Sprintf("Deleted %s", removed)
			}
			return m, nil

		case "t":
			m.CurrentTheme = (m.CurrentTheme + 1) % len(Themes)
			m.FlashNotice = fmt.Sprintf("Theme: %s", Themes[m.CurrentTheme].Name)
			return m, nil

		case "s":
			m.AudioAlert = !m.AudioAlert
			if m.AudioAlert {
				m.FlashNotice = "Audio bell alerts: ON 🔔"
			} else {
				m.FlashNotice = "Audio bell alerts: OFF 🔕"
			}
			return m, nil

		case "r":
			for i := range m.Targets {
				m.Targets[i].Checking = true
			}
			m.FlashNotice = "Refreshing all targets..."
			return m, m.checkAllCmd()
		}

	case checkDoneMsg:
		return m.handleCheckDone(msg)
	}

	return m, nil
}

// handleCheckDone updates the probe result and rings bell on outages if enabled.
func (m Model) handleCheckDone(msg checkDoneMsg) (tea.Model, tea.Cmd) {
	if msg.Index >= 0 && msg.Index < len(m.Targets) {
		m.Targets[msg.Index].Result = &msg.Result
		m.Targets[msg.Index].Checking = false

		if !msg.Result.Up && m.AudioAlert {
			m.FlashNotice = fmt.Sprintf("⚠️ Outage detected on %s!", msg.Result.Target.Raw)
			return m, ringBellCmd
		}
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

// checkTargetCmd produces a tea.Cmd that runs checker.Check in the background.
func checkTargetCmd(index int, target checker.Target, timeout time.Duration) tea.Cmd {
	return func() tea.Msg {
		res := checker.Check(context.Background(), target, timeout)
		return checkDoneMsg{
			Index:  index,
			Result: res,
		}
	}
}
