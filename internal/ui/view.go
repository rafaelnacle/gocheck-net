package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the terminal output string based on the current model state.
func (m Model) View() string {
	if m.Quitting {
		return "\n  Exiting gocheck-net. Have a great day! 🌐\n\n"
	}

	theme := Themes[m.CurrentTheme]
	var b strings.Builder

	// 1. Title and Mascot Header
	b.WriteString("\n ")
	b.WriteString(theme.TitleStyle.Render(" GOCHECK-NET "))
	b.WriteString("  ")
	b.WriteString(theme.SubtitleStyle.Render("Interactive Network Playground"))
	b.WriteString("\n")

	// Calculate counts and health percentage
	upCount := 0
	downCount := 0
	pendingCount := 0
	totalCount := len(m.Targets)

	for _, item := range m.Targets {
		if item.Checking {
			pendingCount++
		} else if item.Result != nil {
			if item.Result.Up {
				upCount++
			} else {
				downCount++
			}
		}
	}

	healthPercent := 0
	if totalCount > 0 {
		healthPercent = (upCount * 100) / totalCount
	}

	// 2. ASCII Mascot & Health Score
	var mascotFace, mascotQuote string
	if pendingCount > 0 && upCount == 0 && downCount == 0 {
		mascotFace = "( o_o )"
		mascotQuote = "Probing network endpoints..."
	} else if downCount == 0 && totalCount > 0 {
		mascotFace = "( ^‿^ )"
		mascotQuote = "All systems operational! 100% Health"
	} else if upCount >= downCount && totalCount > 0 {
		mascotFace = "( •_• )"
		mascotQuote = fmt.Sprintf("%d service(s) down. Review details below.", downCount)
	} else if totalCount == 0 {
		mascotFace = "( ˘_˘ )"
		mascotQuote = "No targets configured. Press [+] to add one!"
	} else {
		mascotFace = "( >_< )"
		mascotQuote = "Major network outage detected! Check connectivity."
	}

	healthGauge := RenderHealthBar(healthPercent, theme)
	mascotLine := fmt.Sprintf(" %s %s  %s",
		theme.MascotStyle.Render(mascotFace),
		theme.SubtitleStyle.Render(fmt.Sprintf("%-42s", mascotQuote)),
		healthGauge,
	)
	b.WriteString(mascotLine)
	b.WriteString("\n\n")

	// 3. Flash Notification Banner (if any)
	if m.FlashNotice != "" {
		noticeStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#2D3748")).
			Padding(0, 1)
		b.WriteString(fmt.Sprintf(" %s\n\n", noticeStyle.Render(m.FlashNotice)))
	}

	// 4. Table Header
	cursorHeader := "  "
	targetHeader := theme.HeaderColStyle.Render(fmt.Sprintf("%-28s", "TARGET"))
	protoHeader := theme.HeaderColStyle.Render(fmt.Sprintf("%-8s", "PROTO"))
	statusHeader := theme.HeaderColStyle.Render(fmt.Sprintf("%-12s", "STATUS"))
	speedHeader := theme.HeaderColStyle.Render(fmt.Sprintf("%-25s", "SPEED GAUGE"))
	detailsHeader := theme.HeaderColStyle.Render("DETAILS")

	b.WriteString(fmt.Sprintf(" %s%s %s %s %s %s\n",
		cursorHeader, targetHeader, protoHeader, statusHeader, speedHeader, detailsHeader))

	// 5. Table Rows
	if len(m.Targets) == 0 {
		b.WriteString(fmt.Sprintf("\n   %s\n\n", theme.SubtitleStyle.Render("No targets in list. Press [+] to add a website or host:port!")))
	} else {
		for i, item := range m.Targets {
			// Cursor indicator
			var cursorStr string
			if i == m.Cursor {
				cursorStr = theme.CursorColStyle.Render("❯ ")
			} else {
				cursorStr = "  "
			}

			// Target name truncated
			addr := item.Target.Raw
			if len(addr) > 26 {
				addr = addr[:23] + "..."
			}
			targetCol := theme.TargetColStyle.Render(fmt.Sprintf("%-28s", addr))
			protoCol := theme.ProtocolColStyle.Render(fmt.Sprintf("%-8s", string(item.Target.Protocol)))

			var statusCol, speedCol, detailsCol string

			if item.Checking {
				statusCol = fmt.Sprintf("%-12s", theme.BadgePendingStyle.Render("CHECKING"))
				speedCol = fmt.Sprintf("%-25s", RenderSpeedBar(0, true, false, theme))
				detailsCol = theme.SubtitleStyle.Render("probing...")
			} else if item.Result != nil {
				res := item.Result
				if res.Up {
					statusCol = fmt.Sprintf("%-12s", theme.BadgeUpStyle.Render("   UP   "))
				} else {
					statusCol = fmt.Sprintf("%-12s", theme.BadgeDownStyle.Render("  DOWN  "))
				}

				speedCol = fmt.Sprintf("%-25s", RenderSpeedBar(res.Latency, false, res.Up, theme))

				if res.Err != nil {
					errStr := res.Err.Error()
					if len(errStr) > 35 {
						errStr = errStr[:32] + "..."
					}
					detailsCol = theme.LatencySlowStyle.Render(errStr)
				} else if res.StatusCode > 0 {
					detailsCol = theme.SubtitleStyle.Render(fmt.Sprintf("HTTP %d OK", res.StatusCode))
				} else {
					detailsCol = theme.SubtitleStyle.Render("OK")
				}
			}

			b.WriteString(fmt.Sprintf(" %s%s %s %s %s %s\n",
				cursorStr, targetCol, protoCol, statusCol, speedCol, detailsCol))
		}
	}

	// 6. Interactive Input Box (when [+] or [a] is active)
	if m.InputMode {
		b.WriteString("\n ")
		inputBox := theme.InputBoxStyle.Render(
			fmt.Sprintf("%s\n%s",
				m.TextInput.View(),
				theme.SubtitleStyle.Render(" [Enter] Add Target   [Esc] Cancel"),
			),
		)
		b.WriteString(inputBox)
		b.WriteString("\n")
	}

	// 7. Stats Box
	statsText := fmt.Sprintf("Total: %d  •  Up: %d  •  Down: %d  •  In-flight: %d  •  Health: %d%%",
		totalCount, upCount, downCount, pendingCount, healthPercent)
	b.WriteString("\n ")
	b.WriteString(theme.StatsBoxStyle.Render(statsText))
	b.WriteString("\n\n")

	// 8. Footer Help Bar
	bellStatus := "ON 🔔"
	if !m.AudioAlert {
		bellStatus = "OFF 🔕"
	}

	b.WriteString("  ")
	b.WriteString(theme.HelpKeyStyle.Render("[↑/↓]"))
	b.WriteString(theme.HelpDescStyle.Render(" Select  "))
	b.WriteString(theme.HelpKeyStyle.Render("[+]"))
	b.WriteString(theme.HelpDescStyle.Render(" Add  "))
	b.WriteString(theme.HelpKeyStyle.Render("[x]"))
	b.WriteString(theme.HelpDescStyle.Render(" Delete  "))
	b.WriteString(theme.HelpKeyStyle.Render("[t]"))
	b.WriteString(theme.HelpDescStyle.Render(fmt.Sprintf(" Theme (%s)  ", theme.Name)))
	b.WriteString(theme.HelpKeyStyle.Render("[s]"))
	b.WriteString(theme.HelpDescStyle.Render(fmt.Sprintf(" Bell (%s)  ", bellStatus)))
	b.WriteString(theme.HelpKeyStyle.Render("[r]"))
	b.WriteString(theme.HelpDescStyle.Render(" Refresh  "))
	b.WriteString(theme.HelpKeyStyle.Render("[q]"))
	b.WriteString(theme.HelpDescStyle.Render(" Quit\n"))

	return b.String()
}
