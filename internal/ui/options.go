package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	optionsMenu struct {
		menu
		mainMenu
		Help  help.Model
		Keys  mainMenuKeys
		Style lipgloss.Style
		Quit  bool
	}
)

func OptionsInitialModel(main mainMenu) tea.Model {
	keys := mainMenuKeys{
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

	return optionsMenu{
		menu: menu{
			Title: "Options",
			Options: []option{
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
		Keys:     keys,
		Help:     help.New(),
		Style:    appStyle,
		mainMenu: main,
	}
}

func (m optionsMenu) Init() tea.Cmd {
	return nil
}

func (m optionsMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch {

		case key.Matches(msg, m.Keys.Next):
			if m.menu.Current < len(m.menu.Options)-1 {
				m.menu.Current++
			}
			return m, nil

		case key.Matches(msg, m.Keys.Prev):
			if m.menu.Current > 0 {
				m.menu.Current--
			}
			return m, nil

		case key.Matches(msg, m.Keys.Accept):
			// WIP
			action := m.menu.Options[m.Current].Name
			switch action {
			case "Quit":
				m.Quit = true
				return m, tea.Quit
			case "Back":
				return m.mainMenu, nil
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

func (m optionsMenu) View() string {
	if m.Quit {
		return ""
	}

	s := m.menu.Title + "\n\n"

	for i, opt := range m.menu.Options {
		indicator := " "
		if i == m.Current {
			indicator = "▶"
		}

		s += fmt.Sprintf("%s %s\n\n", indicator, opt.Name)
	}

	s += m.Help.View(m.Keys)

	return m.Style.Render(s)
}
