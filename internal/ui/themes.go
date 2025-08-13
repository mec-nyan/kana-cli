package ui

import (
	. "github.com/mec-nyan/kana-cli/internal/palette"

	"github.com/charmbracelet/lipgloss"
)

var (
	defaultAppStyle = theme{
		global: lipgloss.NewStyle().
			Padding(1, padding).
			Foreground(lipgloss.ANSIColor(2)),

		input: lipgloss.NewStyle().
			Bold(true),

		highlight: lipgloss.NewStyle().
			Bold(true),

		hints:     lipgloss.NewStyle().
			Italic(true),
	}

	mochaStyle = theme{
		global: lipgloss.NewStyle().
			Padding(1, padding).
			Foreground(lipgloss.Color(Lavender)),

		input: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Teal)),

		highlight: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(Mauve)),

		hints: lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color(Overlay0)),
	}
)
