package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	mainMenu struct {
		menu
		Help help.Model
		Keys mainMenuKeys
		gameOptions
		Quit bool
	}
)

func (k mainMenuKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Show}
}

func (k mainMenuKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Show, k.Next, k.Prev, k.Accept, k.Quit},
	}
}

func MainMenuInitialModel(opts CLIOptions) tea.Model {
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

	return mainMenu{
		menu: menu{
			Title: "Welcome!",
			Options: []option{
				{
					Name: "Start",
				},
				{
					Name: "Options",
				},
				{
					Name: "Help",
				},
				{
					Name: "Quit",
				},
			},
		},
		Keys: keys,
		Help: help.New(),
		gameOptions: gameOptions{
			CLIOptions: opts,
			theme:      defaultAppStyle,
		},
	}
}

func (m mainMenu) Init() tea.Cmd {
	return nil
}

func (m mainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.theme.global = m.theme.global.Width(msg.Width)
		return m, nil

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
			case "Start":
				tm := GameInitialModel(m, m.gameOptions)
				km, _ := tm.(gameModel)
				km.Progress.Width = min(m.theme.global.GetWidth()-padding*2, maxBarWidth)
				return km, nil
			case "Quit":
				m.Quit = true
				return m, tea.Quit
			case "Options":
				return OptionsInitialModel(m), nil
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

func (m mainMenu) View() string {
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
