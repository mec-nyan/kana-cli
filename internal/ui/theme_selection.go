package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	themeMenu struct {
		menu
		optionsMenu
		Help help.Model
		Keys colorMenuKeys
		Quit bool
	}
)

func (k colorMenuKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Show}
}

func (k colorMenuKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Show, k.Next, k.Prev, k.Preview, k.Accept, k.Quit},
	}
}

func ThemeMenuInitialModel(prev optionsMenu) tea.Model {
	keys := colorMenuKeys{
		mainMenuKeys: mainMenuKeys{
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
				key.WithHelp("enter", "apply"),
			),
			Quit: key.NewBinding(
				key.WithKeys(tea.KeyCtrlC.String(), tea.KeyEsc.String(), "q"),
				key.WithHelp("q", "quit"),
			),
		},
		Preview: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "preview"),
		),
	}

	return themeMenu{
		menu: menu{
			Title: "Colours",
			Options: []option{
				{
					Name: "Default",
				},
				{
					Name: "Mocha",
				},
				{
					Name: "Blue",
				},
				{
					Name: "Back",
				},
			},
		},
		Keys:        keys,
		Help:        help.New(),
		optionsMenu: prev,
	}
}

func (m themeMenu) Init() tea.Cmd {
	return nil
}

func (m themeMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlL.String() {
			return m, tea.ClearScreen
		}

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

		case key.Matches(msg, m.Keys.Preview):
			action := m.menu.Options[m.Current].Name
			switch action {
			case "Default":
				m.theme = defaultAppStyle
				return m, nil
			case "Mocha":
				m.theme = mochaStyle
				return m, nil
			case "Purple":
				m.theme = blueStyle
				return m, nil
			default:
				return m, nil
			}

		case key.Matches(msg, m.Keys.Accept):
			// WIP
			action := m.menu.Options[m.Current].Name
			switch action {
			// TODO: Find a better way to handle themes and themes variables.
			case "Default":
				m.theme = defaultAppStyle
				return m.optionsMenu, nil
			case "Mocha":
				m.theme = mochaStyle
				return m.optionsMenu, nil
			case "Purple":
				m.theme = blueStyle
				return m.optionsMenu, nil
			case "Quit":
				m.Quit = true
				return m, tea.Quit
			case "Back":
				return m.optionsMenu, nil
			default:
				return m.optionsMenu, nil
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

func (m themeMenu) View() string {
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

	return m.theme.global.Render(s)
}
