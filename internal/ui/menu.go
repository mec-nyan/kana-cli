package ui

import (
	"fmt"

	. "github.com/mec-nyan/kana-cli/internal/palette"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const padding = 4

var appStyle = lipgloss.NewStyle().Padding(1, padding).Foreground(lipgloss.Color(Lavender))

type (
	Option struct {
		Name     string
		Selected bool
	}

	Menu struct {
		Title   string
		Options []Option
		Current int
	}

	keyMap struct {
		Next   key.Binding
		Prev   key.Binding
		Accept key.Binding
		Quit   key.Binding
	}

	Model struct {
		Menu
		Help  help.Model
		Keys  keyMap
		Style lipgloss.Style
		Quit  bool
	}
)

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Accept, k.Quit},
	}
}

func InitialModel() tea.Model {
	keys := keyMap{
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

	return Model{
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
		Keys: keys,
		Help: help.New(),
		Style: appStyle,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
			// TODO:
			return m, nil
		case key.Matches(msg, m.Keys.Quit):
			m.Quit = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quit {
		return m.Style.Render("Bye!")
	}

	s := m.Menu.Title + "\n\n"

	for i, opt := range m.Menu.Options {
		indicator := " "
		if i == m.Current {
			indicator = ">"
		}

		s += fmt.Sprintf("%s %s\n\n", indicator, opt.Name)
	}

	s += m.Help.View(m.Keys)

	return m.Style.Render(s)
}
