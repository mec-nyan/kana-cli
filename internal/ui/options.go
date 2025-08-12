package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	OptionsModel struct {
		Menu
		MainMenu MainMenuModel
		Help     help.Model
		Keys     keyMap
		Style    lipgloss.Style
		Quit     bool
	}
)

func OptionsInitialModel(main MainMenuModel) tea.Model {
	keys := keyMap{
		Show: key.NewBinding(
			key.WithKeys(";"),
			key.WithHelp(";", "toggle keys"),
		),
		Next: key.NewBinding(
			key.WithKeys(tea.KeyCtrlN.String(), "j"),
			key.WithHelp("j", "next"),
		),
		Prev: key.NewBinding(
			key.WithKeys(tea.KeyCtrlP.String(), "k"),
			key.WithHelp("k", "prev"),
		),
		Accept: key.NewBinding(
			key.WithKeys(tea.KeyEnter.String(), " "),
			key.WithHelp("enter", "accept"),
		),
		Quit: key.NewBinding(
			key.WithKeys(tea.KeyCtrlC.String(), tea.KeyEsc.String(), "q"),
			key.WithHelp("q", "quit"),
		),
	}

	return OptionsModel{
		Menu: Menu{
			Title: "Options",
			Options: []Option{
				{
					Name: "Sound",
				},
				{
					Name: "Syllabary",
				},
				{
					Name: "Back",
				},
				{
					Name: "Help",
				},
				{
					Name: "Quit",
				},
			},
		},
		Keys:  keys,
		Help:  help.New(),
		Style: appStyle,
		MainMenu: main,
	}
}

func (m OptionsModel) Init() tea.Cmd {
	return nil
}

func (m OptionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch {

		case key.Matches(msg, m.Keys.Next):
			if m.Menu.Current < len(m.Menu.Options)-1 {
				m.Menu.Current++
			}
			return m, nil

		case key.Matches(msg, m.Keys.Prev):
			if m.Menu.Current > 0 {
				m.Menu.Current--
			}
			return m, nil

		case key.Matches(msg, m.Keys.Accept):
			// WIP
			action := m.Menu.Options[m.Current].Name
			switch action {
			case "Quit":
				m.Quit = true
				return m, tea.Quit
			case "Back":
				return m.MainMenu, nil
			default:
				return m, nil
			}

		case key.Matches(msg, m.Keys.Quit):
			m.Quit = true
			return m, tea.Quit

		case key.Matches(msg, m.Keys.Show):
			m.Help.ShowAll = !m.Help.ShowAll
			return m, nil
		}
	}

	return m, nil
}

func (m OptionsModel) View() string {
	if m.Quit {
		return ""
	}

	s := m.Menu.Title + "\n\n"

	for i, opt := range m.Menu.Options {
		indicator := " "
		if i == m.Current {
			indicator = "▶"
		}

		s += fmt.Sprintf("%s %s\n\n", indicator, opt.Name)
	}

	s += m.Help.View(m.Keys)

	return m.Style.Render(s)
}
