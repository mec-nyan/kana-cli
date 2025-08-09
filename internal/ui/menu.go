package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type option struct {
	Name     string
	Selected bool
}

type step struct {
	Name    string
	Options []option
}

type Model struct {
	Cursor      int
	Steps       []step
	CurrentStep int
}

func InitialModel() tea.Model {
	return Model{
		Cursor: 0,
		CurrentStep: 0,
		Steps: []step{
			{
				Name: "Start",
				Options: []option{
					{
						Name: "new game",
					},
					{
						Name: "saved",
					},
				},
			},
			{
				Name: "Syllabary",
				Options: []option{
					{
						Name: "hiragana",
					},
					{
						Name: "katakana",
					},
				},
			},
			{
				Name: "Sound",
				Options: []option{
					{
						Name: "on",
					},
					{
						Name: "off",
					},
				},
			},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	options := m.Steps[m.CurrentStep].Options

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+n", "j":
			if m.Cursor < len(options)-1 {
				m.Cursor++
			}

		case "ctrl+p", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "enter", " ":
			for i := range options {
				if i == m.Cursor {
					options[i].Selected = true
				} else {
					options[i].Selected = false
				}
			}
			m.Steps[m.CurrentStep].Options = options
			m.CurrentStep++
			m.Cursor = 0

			// For now, just quit.
			if m.CurrentStep == len(m.Steps) {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.CurrentStep == len(m.Steps) {
		return ""
	}

	options := m.Steps[m.CurrentStep]

	s := options.Name + ":\n\n"

	for i, choice := range options.Options {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("  %s %s\n", cursor, choice.Name)
	}

	s += "\n\n[j] next - [k] prev - [q] quit\n"

	return s
}
