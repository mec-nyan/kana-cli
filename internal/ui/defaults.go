package ui

import (
	. "github.com/mec-nyan/kana-cli/internal/palette"

	"github.com/charmbracelet/lipgloss"
)

const padding = 4

var appStyle = lipgloss.NewStyle().Padding(1, padding).Foreground(lipgloss.Color(Lavender))
