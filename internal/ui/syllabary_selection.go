package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	syllabaryMenu struct {
		menu
		optionsMenu
		Help help.Model
		Keys mainMenuKeys
		Quit bool
	}
)


func SyllabaryMenuInitialModel(prev optionsMenu) tea.Model {
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

	return syllabaryMenu{
		menu: menu{
			Title: "Practise",
			Options: []option{
				{
					Name: "Hiragana",
				},
				{
					Name: "Katakana",
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

func (m syllabaryMenu) Init() tea.Cmd {
	return nil
}

func (m syllabaryMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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

		case key.Matches(msg, m.Keys.Accept):
			action := m.menu.Options[m.Current].Name
			switch action {
			case "Hiragana":
				m.syllabary = hiragana
				return m.optionsMenu, nil
			case "Katakana":
				m.syllabary = katakana
				return m.optionsMenu, nil
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

func (m syllabaryMenu) View() string {
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
