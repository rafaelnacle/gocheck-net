package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Theme encapsulates the Lip Gloss styling for a color scheme.
type Theme struct {
	Name              string
	TitleStyle        lipgloss.Style
	SubtitleStyle     lipgloss.Style
	HeaderColStyle    lipgloss.Style
	TargetColStyle    lipgloss.Style
	CursorColStyle    lipgloss.Style
	ProtocolColStyle  lipgloss.Style
	BadgeUpStyle      lipgloss.Style
	BadgeDownStyle    lipgloss.Style
	BadgePendingStyle lipgloss.Style
	LatencyFastStyle  lipgloss.Style
	LatencyMedStyle   lipgloss.Style
	LatencySlowStyle  lipgloss.Style
	StatsBoxStyle     lipgloss.Style
	HelpKeyStyle      lipgloss.Style
	HelpDescStyle     lipgloss.Style
	MascotStyle       lipgloss.Style
	InputBoxStyle     lipgloss.Style
}

// Themes holds the list of selectable palettes.
var Themes = []Theme{
	// 0: Modern Default (Indigo & Emerald)
	{
		Name: "Default Modern",
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5A56E0")).
			Padding(0, 1),
		SubtitleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true),
		HeaderColStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A0A0A0")).
			Underline(true),
		TargetColStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")),
		CursorColStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04D182")),
		ProtocolColStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D8296")),
		BadgeUpStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#04D182")).
			Padding(0, 1),
		BadgeDownStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF3355")).
			Padding(0, 1),
		BadgePendingStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFB800")).
			Padding(0, 1),
		LatencyFastStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04D182")),
		LatencyMedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB800")),
		LatencySlowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF3355")),
		StatsBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3C3C48")).
			Padding(0, 1),
		HelpKeyStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")),
		HelpDescStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")),
		MascotStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00E5FF")),
		InputBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5A56E0")).
			Padding(0, 1),
	},
	// 1: Matrix Hacker (CRT Phosphor Green)
	{
		Name: "Matrix Hacker",
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#00FF66")).
			Padding(0, 1),
		SubtitleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#008833")).
			Italic(true),
		HeaderColStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00CC55")).
			Underline(true),
		TargetColStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#88FFAA")),
		CursorColStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF66")),
		ProtocolColStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#009944")),
		BadgeUpStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#00FF66")).
			Padding(0, 1),
		BadgeDownStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FF2255")).
			Padding(0, 1),
		BadgePendingStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#DDFF00")).
			Padding(0, 1),
		LatencyFastStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF66")),
		LatencyMedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDFF00")),
		LatencySlowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4466")),
		StatsBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00AA44")).
			Padding(0, 1),
		HelpKeyStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF66")),
		HelpDescStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#008833")),
		MascotStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF66")),
		InputBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FF66")).
			Padding(0, 1),
	},
	// 2: Cyberpunk Neon (Hot Pink, Cyan, Violet)
	{
		Name: "Cyberpunk Neon",
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF007F")).
			Padding(0, 1),
		SubtitleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00F0FF")).
			Italic(true),
		HeaderColStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF77AA")).
			Underline(true),
		TargetColStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0F7FA")),
		CursorColStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00F0FF")),
		ProtocolColStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AA00FF")),
		BadgeUpStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#00F0FF")).
			Padding(0, 1),
		BadgeDownStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF0055")).
			Padding(0, 1),
		BadgePendingStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFE600")).
			Padding(0, 1),
		LatencyFastStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00F0FF")),
		LatencyMedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE600")),
		LatencySlowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0055")),
		StatsBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF007F")).
			Padding(0, 1),
		HelpKeyStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF007F")),
		HelpDescStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8888AA")),
		MascotStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF007F")),
		InputBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00F0FF")).
			Padding(0, 1),
	},
	// 3: Retro Amber (1980s Vintage Terminal)
	{
		Name: "Retro Amber",
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFB000")).
			Padding(0, 1),
		SubtitleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C68400")).
			Italic(true),
		HeaderColStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFA000")).
			Underline(true),
		TargetColStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE099")),
		CursorColStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFB000")),
		ProtocolColStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AA7000")),
		BadgeUpStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFB000")).
			Padding(0, 1),
		BadgeDownStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#CC4400")).
			Padding(0, 1),
		BadgePendingStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFD54F")).
			Padding(0, 1),
		LatencyFastStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB000")),
		LatencyMedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA000")),
		LatencySlowStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CC4400")),
		StatsBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#AA7000")).
			Padding(0, 1),
		HelpKeyStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFB000")),
		HelpDescStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#885500")),
		MascotStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD54F")),
		InputBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFB000")).
			Padding(0, 1),
	},
}

// RenderSpeedBar draws a 10-block visual bar with latency and speed emoji.
func RenderSpeedBar(latency time.Duration, checking bool, isUp bool, theme Theme) string {
	const barWidth = 10

	if checking {
		return theme.SubtitleStyle.Render("[░░░░░░░░░░] ⏳ probing")
	}

	if !isUp {
		return theme.LatencySlowStyle.Render("[░░░░░░░░░░] 💥 failed")
	}

	// Calculate filled blocks: lower latency = more filled / higher speed rating
	// Range: 0ms (10 blocks) to 1000ms (1 block)
	ms := latency.Milliseconds()
	var filled int
	var emoji string
	var style lipgloss.Style

	switch {
	case ms < 80:
		filled = 10
		emoji = "⚡"
		style = theme.LatencyFastStyle
	case ms < 150:
		filled = 8
		emoji = "🚀"
		style = theme.LatencyFastStyle
	case ms < 300:
		filled = 6
		emoji = "✨"
		style = theme.LatencyFastStyle
	case ms < 600:
		filled = 4
		emoji = "🚶"
		style = theme.LatencyMedStyle
	default:
		filled = 2
		emoji = "🐌"
		style = theme.LatencySlowStyle
	}

	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	bar := fmt.Sprintf("[%s%s] %s %-6s",
		strings.Repeat("█", filled),
		strings.Repeat("░", empty),
		emoji,
		latency.Round(time.Millisecond).String(),
	)

	return style.Render(bar)
}

// RenderHealthBar draws a visual health gauge for total uptime percentage.
func RenderHealthBar(percent int, theme Theme) string {
	const barWidth = 14
	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		percent = 100
	}

	filled := (percent * barWidth) / 100
	empty := barWidth - filled

	var style lipgloss.Style
	if percent >= 80 {
		style = theme.LatencyFastStyle
	} else if percent >= 50 {
		style = theme.LatencyMedStyle
	} else {
		style = theme.LatencySlowStyle
	}

	gauge := fmt.Sprintf("[%s%s] %d%%",
		strings.Repeat("█", filled),
		strings.Repeat("░", empty),
		percent,
	)

	return style.Render(gauge)
}
