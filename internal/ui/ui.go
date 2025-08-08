package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepMainMenu step = iota
	stepModeMenu
	stepExit
)

type model struct {
	cursor int
	choices []string
	step step
}


func InitialModel() tea.Model {
	return model{
		cursor: 0,
		choices: []string{"new game", "saved"},
		step: stepMainMenu,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
