package ui

import (
	. "github.com/mec-nyan/kana-cli/internal/palette"

	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, padding).
			Foreground(lipgloss.Color(Lavender))

	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(Mauve))

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Teal))

	hintStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color(Overlay0))
)
