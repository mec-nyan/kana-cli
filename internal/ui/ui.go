package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepMainMenu step = iota
	stepModeMenu
	stepExit
)

type option struct {
	name string
	selected bool
}

type model struct {
	cursor  int
	choices []option
	step    step
}

func InitialModel() tea.Model {
	return model{
		cursor:  0,
		choices: []option{
			{
				name: "new game",
				selected: false,
			},
			{
				name: "saved",
				selected: false,
			},
		},
		step:    stepMainMenu,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+n", "j":
			if m.cursor < len(m.choices) - 1 {
				m.cursor++
			}

		case "ctrl+p", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter", " ":
			for i := range m.choices {
				if i == m.cursor {
					m.choices[i].selected = !m.choices[i].selected
				} else {
					m.choices[i].selected = false
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Select an option:\n"

	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "->"
		}

		checked := " "
		if choice.selected {
			checked = "x"
		}

		s += fmt.Sprintf("  %s [%s] %s\n", cursor, checked, choice.name)
	}

	s += "\n\n[j] next - [k] prev - [q] quit\n"

	return s
}
