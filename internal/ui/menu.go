package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	MainMenu struct {
		Menu
		Help  help.Model
		Keys  MenuKeys
		Style lipgloss.Style
		Opts  GameOptions
		Quit  bool
	}
)

func (k MenuKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Show}
}

func (k MenuKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Show, k.Next, k.Prev, k.Accept, k.Quit},
	}
}

func MainMenuInitialModel(opts GameOptions) tea.Model {
	keys := MenuKeys{
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

	return MainMenu{
		Menu: Menu{
			Title: "Welcome!",
			Options: []Option{
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
		Keys:  keys,
		Help:  help.New(),
		Style: appStyle,
		Opts:  opts,
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Style = m.Style.Width(msg.Width)
		return m, nil

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
			case "Start":
				tm := GameInitialModel(m, m.Opts)
				km, _ := tm.(GameModel)
				km.Progress.Width = min(m.Style.GetWidth()-padding*2, maxBarWidth)
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

func (m MainMenu) View() string {
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
