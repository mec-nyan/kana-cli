package ui

import (
	. "github.com/mec-nyan/kana-cli/internal/palette"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

var (
	greenStyle = theme{
		global: lipgloss.NewStyle().
			Padding(1, padding).
			Foreground(lipgloss.ANSIColor(2)),

		input: lipgloss.NewStyle().
			Bold(true),

		highlight: lipgloss.NewStyle().
			Foreground(lipgloss.ANSIColor(10)).
			Bold(true),

		hints: lipgloss.NewStyle().
			Italic(true),

		progressOpts: progressOpts{
			options: []progress.Option{
				progress.WithFillCharacters('▁', '▁'),
			},
			empty:     "60",
			full:      "2",
			showPerc:  true,
			percStyle: lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)),
		},
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

		progressOpts: progressOpts{
			options: []progress.Option{
				progress.WithGradient(Mauve, Sapphire),
				progress.WithFillCharacters('▂', '▂'),
			},
			empty:     Surface0,
			showPerc:  true,
			percStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(Subtext0)),
		},
	}

	blueStyle = theme{
		global: lipgloss.NewStyle().
			Padding(1, padding).
			Foreground(lipgloss.Color(Blue)),

		input: lipgloss.NewStyle().
			Bold(true),

		highlight: lipgloss.NewStyle().
			Foreground(lipgloss.Color(Blue)).
			Bold(true),

		hints: lipgloss.NewStyle().
			Italic(true),

		progressOpts: progressOpts{
			options: []progress.Option{
				progress.WithGradient(Blue, Lavender),
				progress.WithFillCharacters('▁', '▁'),
			},
			empty:     Surface0,
			showPerc:  true,
			percStyle: lipgloss.NewStyle().Foreground(lipgloss.Color(Blue)),
		},
	}

	defaultAppStyle = mochaStyle
)
