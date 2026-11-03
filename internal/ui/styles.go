package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Header & Title
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5A56E0")).
			Padding(0, 1).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)

	// Table Headers & Rows
	headerColStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A0A0A0")).
			Underline(true)

	targetColStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0"))

	protocolColStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D8296"))

	// Status Badges
	badgeUpStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#04D182")).
			Padding(0, 1)

	badgeDownStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF3355")).
			Padding(0, 1)

	badgePendingStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#FFB800")).
				Padding(0, 1)

	// Latency styling
	latencyFastStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04D182"))

	latencyMedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB800"))

	latencySlowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF3355"))

	// Footer & Help
	footerHelpKeyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4"))

	footerHelpDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#777777"))

	statsBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3C3C48")).
			Padding(0, 1).
			MarginTop(1)
)
