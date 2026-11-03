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

	// If Inspector modal is active, render detailed target overlay
	if m.InspectorMode && len(m.Targets) > 0 && m.Cursor < len(m.Targets) {
		return m.renderInspector(theme)
	}

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
			Background(lipgloss.Color("#2E78D2")).
			Padding(0, 1).
			MarginLeft(1).
			MarginBottom(1)
		b.WriteString(noticeStyle.Render("ℹ "+m.FlashNotice) + "\n\n")
	}

	// 4. Table Headers
	cursorHeader := "  "
	targetHeader := theme.HeaderColStyle.Render(fmt.Sprintf("%-28s", "TARGET"))
	protoHeader := theme.HeaderColStyle.Render(fmt.Sprintf("%-8s", "PROTO"))
	statusHeader := theme.HeaderColStyle.Render(fmt.Sprintf("%-12s", "STATUS"))
	speedHeader := theme.HeaderColStyle.Render(fmt.Sprintf("%-25s", "SPEED GAUGE"))
	detailsHeader := theme.HeaderColStyle.Render("DETAILS")

	b.WriteString(fmt.Sprintf(" %s%s %s %s %s %s\n",
		cursorHeader, targetHeader, protoHeader, statusHeader, speedHeader, detailsHeader))

	// 5. Rows
	for i, item := range m.Targets {
		// Cursor marker
		marker := "  "
		if i == m.Cursor {
			marker = theme.CursorColStyle.Render("❯ ")
		}

		// Truncate target address if long
		addr := item.Target.Raw
		if len(addr) > 26 {
			addr = addr[:23] + "..."
		}
		targetCol := theme.TargetColStyle.Render(fmt.Sprintf("%-28s", addr))
		protoCol := theme.ProtocolColStyle.Render(fmt.Sprintf("%-8s", string(item.Target.Protocol)))

		var statusCol, detailsCol string
		speedCol := RenderSpeedBar(0, item.Checking, false, theme)

		if item.Checking {
			statusCol = fmt.Sprintf("%-12s", theme.BadgePendingStyle.Render("CHECKING"))
			detailsCol = theme.SubtitleStyle.Render("probing...")
		} else if item.Result != nil {
			res := item.Result
			speedCol = RenderSpeedBar(res.Latency, false, res.Up, theme)

			if res.Up {
				statusCol = fmt.Sprintf("%-12s", theme.BadgeUpStyle.Render("   UP   "))
			} else {
				statusCol = fmt.Sprintf("%-12s", theme.BadgeDownStyle.Render("  DOWN  "))
			}

			// Details or error
			if res.Err != nil {
				errStr := res.Err.Error()
				if len(errStr) > 35 {
					errStr = errStr[:32] + "..."
				}
				detailsCol = theme.LatencySlowStyle.Render(errStr)
			} else if item.Target.Protocol == "HTTP" && res.StatusCode > 0 {
				detailsCol = theme.SubtitleStyle.Render(fmt.Sprintf("HTTP %d OK", res.StatusCode))
			} else {
				detailsCol = theme.SubtitleStyle.Render("OK")
			}
		}

		b.WriteString(fmt.Sprintf(" %s%s %s %s %s %s\n",
			marker, targetCol, protoCol, statusCol, speedCol, detailsCol))
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
	b.WriteString(theme.HelpKeyStyle.Render("[Enter]"))
	b.WriteString(theme.HelpDescStyle.Render(" Inspect  "))
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

// renderInspector renders the full-screen modal card for the selected target.
func (m Model) renderInspector(theme Theme) string {
	item := m.Targets[m.Cursor]
	var b strings.Builder

	b.WriteString("\n ")
	b.WriteString(theme.TitleStyle.Render(" TARGET INSPECTOR "))
	b.WriteString("  ")
	b.WriteString(theme.SubtitleStyle.Render(fmt.Sprintf("Target [%d/%d]", m.Cursor+1, len(m.Targets))))
	b.WriteString("\n\n")

	var card strings.Builder

	// Target overview
	card.WriteString(fmt.Sprintf(" %s %s\n",
		theme.HeaderColStyle.Render("Target Endpoint:"),
		theme.CursorColStyle.Render(item.Target.Raw),
	))
	card.WriteString(fmt.Sprintf(" %s %s    %s %s\n",
		theme.SubtitleStyle.Render("Protocol:"),
		theme.TargetColStyle.Render(string(item.Target.Protocol)),
		theme.SubtitleStyle.Render("Address:"),
		theme.TargetColStyle.Render(item.Target.Address),
	))

	// Status & Latency
	if item.Checking {
		card.WriteString(fmt.Sprintf(" %s %s\n\n",
			theme.SubtitleStyle.Render("Status:"),
			theme.BadgePendingStyle.Render("CHECKING IN-FLIGHT"),
		))
	} else if item.Result != nil {
		res := item.Result
		statusBadge := theme.BadgeUpStyle.Render("  UP (ONLINE)  ")
		if !res.Up {
			statusBadge = theme.BadgeDownStyle.Render("  DOWN (OFFLINE)  ")
		}

		codeStr := "-"
		if res.StatusCode > 0 {
			codeStr = fmt.Sprintf("%d", res.StatusCode)
		}

		card.WriteString(fmt.Sprintf(" %s %s    %s %s    %s %s\n\n",
			theme.SubtitleStyle.Render("Status:"),
			statusBadge,
			theme.SubtitleStyle.Render("Status Code:"),
			theme.TargetColStyle.Render(codeStr),
			theme.SubtitleStyle.Render("Latency:"),
			RenderSpeedBar(res.Latency, false, res.Up, theme),
		))

		if res.Err != nil {
			card.WriteString(fmt.Sprintf(" %s\n   %s\n\n",
				theme.LatencySlowStyle.Render("⚠️ Error Details:"),
				theme.SubtitleStyle.Render(res.Err.Error()),
			))
		}

		// DNS Resolution section
		card.WriteString(fmt.Sprintf(" %s\n", theme.HeaderColStyle.Render("── DNS RESOLUTION ─────────────────────────────────────────")))
		if len(res.ResolvedIPs) > 0 {
			card.WriteString(fmt.Sprintf("   Resolved IPs: %s\n\n",
				theme.TargetColStyle.Render(strings.Join(res.ResolvedIPs, ", ")),
			))
		} else {
			card.WriteString("   Resolved IPs: (No IP records resolved)\n\n")
		}

		// SSL/TLS Certificate section
		card.WriteString(fmt.Sprintf(" %s\n", theme.HeaderColStyle.Render("── SSL/TLS CERTIFICATE ────────────────────────────────────")))
		if res.Cert != nil {
			cert := res.Cert
			var certStatusBadge string
			if cert.Expired {
				certStatusBadge = theme.BadgeDownStyle.Render(" EXPIRED ❌ ")
			} else if cert.DaysLeft <= 14 {
				certStatusBadge = theme.BadgePendingStyle.Render(fmt.Sprintf(" EXPIRING SOON: %d DAYS ⚠️ ", cert.DaysLeft))
			} else {
				certStatusBadge = theme.BadgeUpStyle.Render(fmt.Sprintf(" VALID (%d DAYS LEFT) ✅ ", cert.DaysLeft))
			}

			card.WriteString(fmt.Sprintf("   Issuer:       %s\n", theme.TargetColStyle.Render(cert.Issuer)))
			card.WriteString(fmt.Sprintf("   Common Name:  %s\n", theme.TargetColStyle.Render(cert.CommonName)))
			card.WriteString(fmt.Sprintf("   TLS Version:  %s\n", theme.TargetColStyle.Render(cert.TLSVersion)))
			card.WriteString(fmt.Sprintf("   Valid From:   %s\n", theme.SubtitleStyle.Render(cert.NotBefore.Format("2006-01-02 15:04:05 UTC"))))
			card.WriteString(fmt.Sprintf("   Expires:      %s  %s\n",
				theme.TargetColStyle.Render(cert.NotAfter.Format("2006-01-02 15:04:05 UTC")),
				certStatusBadge,
			))
			if len(cert.DNSNames) > 0 {
				sans := strings.Join(cert.DNSNames, ", ")
				if len(sans) > 55 {
					sans = sans[:52] + "..."
				}
				card.WriteString(fmt.Sprintf("   Alt Names:    %s\n", theme.SubtitleStyle.Render(sans)))
			}
			card.WriteString("\n")
		} else {
			card.WriteString("   SSL/TLS: Not applicable (Plain HTTP or non-TLS TCP port)\n\n")
		}

		// HTTP Headers section
		if len(res.Headers) > 0 {
			card.WriteString(fmt.Sprintf(" %s\n", theme.HeaderColStyle.Render("── HTTP RESPONSE HEADERS ──────────────────────────────────")))
			for k, v := range res.Headers {
				card.WriteString(fmt.Sprintf("   %-15s %s\n",
					theme.SubtitleStyle.Render(k+":"),
					theme.TargetColStyle.Render(v),
				))
			}
			card.WriteString("\n")
		}
	}

	// Wrap in a stylized card box
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5A56E0")).
		Padding(1, 2).
		MarginBottom(1)

	b.WriteString(cardStyle.Render(card.String()))
	b.WriteString("\n\n")

	// Modal Navigation Footer
	b.WriteString("  ")
	b.WriteString(theme.HelpKeyStyle.Render("[↑/↓] or [k/j]"))
	b.WriteString(theme.HelpDescStyle.Render(" Switch Target  "))
	b.WriteString(theme.HelpKeyStyle.Render("[Esc/Enter/q]"))
	b.WriteString(theme.HelpDescStyle.Render(" Close Inspector\n"))

	return b.String()
}
