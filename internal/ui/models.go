package ui

import "github.com/charmbracelet/bubbles/key"

type (
	GameOptions struct {
		Test bool
		Auto bool
	}

	Question struct {
		hiragana string
		romaji   []string
		hints    []string
		played   bool
	}

	Option struct {
		Name string
	}

	Menu struct {
		Title   string
		Options []Option
		Current int
	}

	MenuKeys struct {
		Show   key.Binding
		Next   key.Binding
		Prev   key.Binding
		Accept key.Binding
		Quit   key.Binding
	}

	errMsg error

	tickMsg struct{}
)
