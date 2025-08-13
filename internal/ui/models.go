package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type (
	CLIOptions struct {
		Test bool
		Auto bool
	}

	gameOptions struct {
		CLIOptions
		theme
	}

	question struct {
		hiragana string
		romaji   []string
		hints    []string
		played   bool
	}

	option struct {
		Name string
	}

	menu struct {
		Title   string
		Options []option
		Current int
	}

	mainMenuKeys struct {
		Show   key.Binding
		Next   key.Binding
		Prev   key.Binding
		Accept key.Binding
		Quit   key.Binding
	}

	errMsg error

	tickMsg struct{}

	theme struct {
		global    lipgloss.Style
		input     lipgloss.Style
		highlight lipgloss.Style
		hints     lipgloss.Style
	}
)
