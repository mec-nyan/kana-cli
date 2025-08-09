package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type option struct {
	name     string
	selected bool
}

type step struct {
	name    string
	options []option
}

type model struct {
	cursor      int
	steps       []step
	currentStep int
}

func InitialModel() tea.Model {
	return model{
		cursor: 0,
		currentStep: 0,
		steps: []step{
			{
				name: "Start",
				options: []option{
					{
						name: "new game",
					},
					{
						name: "saved",
					},
				},
			},
			{
				name: "Syllabary",
				options: []option{
					{
						name: "hiragana",
					},
					{
						name: "katakana",
					},
				},
			},
			{
				name: "Play sound",
				options: []option{
					{
						name: "on",
					},
					{
						name: "off",
					},
				},
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	options := m.steps[m.currentStep].options

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+n", "j":
			if m.cursor < len(options)-1 {
				m.cursor++
			}

		case "ctrl+p", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter", " ":
			for i := range options {
				if i == m.cursor {
					options[i].selected = true
				} else {
					options[i].selected = false
				}
			}
			m.steps[m.currentStep].options = options
			m.currentStep++
			m.cursor = 0

			// For now, just quit.
			if m.currentStep == len(m.steps) {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.currentStep == len(m.steps) {
		return "END"
	}

	options := m.steps[m.currentStep]

	s := options.name + ":\n\n"

	for i, choice := range options.options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("  %s %s\n", cursor, choice.name)
	}

	s += "\n\n[j] next - [k] prev - [q] quit\n"

	return s
}
