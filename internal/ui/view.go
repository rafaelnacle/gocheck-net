package ui

import (
	"fmt"
	"strings"
	"time"
)

// View renders the terminal output string based on current state.
func (m Model) View() string {
	if m.Quitting {
		return "\n  Exiting gocheck-net. Goodbye!\n\n"
	}

	var b strings.Builder

	// Title header
	b.WriteString("\n ")
	b.WriteString(titleStyle.Render(" GOCHECK-NET "))
	b.WriteString("  ")
	b.WriteString(subtitleStyle.Render("Interactive Network Status Checker"))
	b.WriteString("\n\n")

	// Table Headers
	targetHeader := headerColStyle.Render(fmt.Sprintf("%-34s", "TARGET"))
	protoHeader := headerColStyle.Render(fmt.Sprintf("%-10s", "PROTOCOL"))
	statusHeader := headerColStyle.Render(fmt.Sprintf("%-12s", "STATUS"))
	codeHeader := headerColStyle.Render(fmt.Sprintf("%-8s", "CODE"))
	latencyHeader := headerColStyle.Render(fmt.Sprintf("%-12s", "LATENCY"))
	detailsHeader := headerColStyle.Render("DETAILS")

	b.WriteString(fmt.Sprintf(" %s %s %s %s %s %s\n",
		targetHeader, protoHeader, statusHeader, codeHeader, latencyHeader, detailsHeader))

	// Rows
	upCount := 0
	downCount := 0
	pendingCount := 0

	for _, item := range m.Targets {
		// Target address truncated if too long
		addr := item.Target.Raw
		if len(addr) > 32 {
			addr = addr[:29] + "..."
		}
		targetCol := targetColStyle.Render(fmt.Sprintf("%-34s", addr))
		protoCol := protocolColStyle.Render(fmt.Sprintf("%-10s", string(item.Target.Protocol)))

		var statusCol, codeCol, latencyCol, detailsCol string

		if item.Checking {
			pendingCount++
			statusCol = fmt.Sprintf("%-12s", badgePendingStyle.Render("CHECKING"))
			codeCol = fmt.Sprintf("%-8s", "-")
			latencyCol = fmt.Sprintf("%-12s", "-")
			detailsCol = subtitleStyle.Render("probing...")
		} else if item.Result != nil {
			res := item.Result
			if res.Up {
				upCount++
				statusCol = fmt.Sprintf("%-12s", badgeUpStyle.Render("   UP   "))
			} else {
				downCount++
				statusCol = fmt.Sprintf("%-12s", badgeDownStyle.Render("  DOWN  "))
			}

			// HTTP status code
			if item.Target.Protocol == "HTTP" && res.StatusCode > 0 {
				codeCol = fmt.Sprintf("%-8d", res.StatusCode)
			} else {
				codeCol = fmt.Sprintf("%-8s", "-")
			}

			// Latency styling
			dur := res.Latency.Round(time.Millisecond)
			durStr := fmt.Sprintf("%-12s", dur.String())
			if res.Latency < 200*time.Millisecond {
				latencyCol = latencyFastStyle.Render(durStr)
			} else if res.Latency < 800*time.Millisecond {
				latencyCol = latencyMedStyle.Render(durStr)
			} else {
				latencyCol = latencySlowStyle.Render(durStr)
			}

			// Details or error
			if res.Err != nil {
				errStr := res.Err.Error()
				if len(errStr) > 40 {
					errStr = errStr[:37] + "..."
				}
				detailsCol = latencySlowStyle.Render(errStr)
			} else {
				detailsCol = subtitleStyle.Render("OK")
			}
		}

		b.WriteString(fmt.Sprintf(" %s %s %s %s %s %s\n",
			targetCol, protoCol, statusCol, codeCol, latencyCol, detailsCol))
	}

	// Statistics Bar
	statsText := fmt.Sprintf("Total: %d  •  Up: %d  •  Down: %d  •  In-flight: %d",
		len(m.Targets), upCount, downCount, pendingCount)
	b.WriteString("\n ")
	b.WriteString(statsBoxStyle.Render(statsText))
	b.WriteString("\n\n")

	// Footer Help
	b.WriteString("  ")
	b.WriteString(footerHelpKeyStyle.Render("[r]"))
	b.WriteString(footerHelpDescStyle.Render(" Refresh  "))
	b.WriteString(footerHelpKeyStyle.Render("[q/esc]"))
	b.WriteString(footerHelpDescStyle.Render(" Quit"))
	b.WriteString("\n")

	return b.String()
}
