package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/mec-nyan/kana-master/pkg/kana"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	gameState struct {
		current  int
		tries    int
		percent  int
		quit     bool
		err      error
		hint     bool
		end      bool
		accuracy float64
	}

	gameModel struct {
		questions []question
		gameState
		gameOptions
		textInput textinput.Model
		Progress  progress.Model
		keys      gameKeys
		help      help.Model
		menu      mainMenu
	}

	gameKeys struct {
		Show    key.Binding
		Accept  key.Binding
		Menu    key.Binding
		Command key.Binding
		Hint    key.Binding
		Help    key.Binding
		Quit    key.Binding
	}
)

func (k gameKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Show}
}

func (k gameKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Show, k.Accept, k.Menu, k.Hint, k.Command, k.Help, k.Quit},
	}
}

func GameInitialModel(menu mainMenu, opts gameOptions) tea.Model {

	questions := makeQuestions(opts)
	ti := makeGameInput(opts.theme.input)
	pBar := makeProgressBar(opts.progressOpts)
	keys := makeGameKeys()

	return gameModel{
		questions:   questions,
		gameOptions: opts,
		textInput:   ti,
		Progress:    pBar,
		keys:        keys,
		help:        help.New(),
		menu:        menu,
	}
}

func (m gameModel) Init() tea.Cmd {
	return nil
}

func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var q question
	if m.current < len(m.questions) {
		q = m.questions[m.current]
	}

	switch msg := msg.(type) {

	case tickMsg:
		m.quit = true
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.theme.global = m.theme.global.Width(msg.Width)
		m.Progress.Width = min(msg.Width-padding*2, maxBarWidth)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlL.String() {
			return m, tea.ClearScreen
		}

		switch {

		// There's no hiragana/katakana that starts with "q".
		// We can safely use this to quit.
		case key.Matches(msg, m.keys.Quit):
			m.quit = true
			return m, tea.Quit

		// Use both enter or space to accept input value.
		case key.Matches(msg, m.keys.Accept):
			return m.checkAnswer(q)

		case key.Matches(msg, m.keys.Show):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, m.keys.Hint):
			m.hint = !m.hint
			return m, nil

		case key.Matches(msg, m.keys.Help):
			// TODO
			return m, nil

		case key.Matches(msg, m.keys.Command):
			// TODO
			return m, nil

		case key.Matches(msg, m.keys.Menu):
			return m.menu, nil

		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m gameModel) View() string {
	if m.quit {
		return "Good bye then!"
	}

	var s string

	if m.end {
		s = fmt.Sprintf(`%s

You've done it!

Accuracy: %0.1f%%

%s`,
			m.Progress.ViewAs(1),
			float64(len(m.questions))/float64(m.tries)*100,
			m.help.View(m.keys))

		return m.theme.global.Render(s)
	}

	if m.current >= len(m.questions) {
		return "Nope!"
	}

	q := m.questions[m.current]

	hint := "Hint: ..."
	if m.hint {
		hint = "Hint: " + strings.Join(q.hints, ", ")
	}

	progress := 100.0 / float64(len(m.questions)) * float64(m.current)
	if m.tries > 0 {
		m.accuracy = float64(m.current) / float64(m.tries) * 100.0
	}

	// TODO: Use lipgloss for padding and styling!
	s = fmt.Sprintf(`%s

Hiragana: %s

Write in romaji: %s

%s 

Accuracy: %0.1f%%

%s`,
		m.Progress.ViewAs(progress/100),
		m.theme.highlight.Render(q.kana),
		m.textInput.View(),
		m.theme.hints.Render(hint),
		m.accuracy,
		m.help.View(m.keys))

	return m.theme.global.Render(s)
}

// TODO: How to check in autoMode?
func (m gameModel) checkAnswer(q question) (tea.Model, tea.Cmd) {
	m.hint = false
	m.tries++
	guess := m.textInput.Value()
	for _, rmj := range q.romaji {
		if guess == rmj {
			m.current++
			if m.current == len(m.questions) {
				m.end = true
				return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
					return tickMsg{}
				})
			}
		}
	}
	m.textInput.Reset()
	return m, nil
}

// Utility functions.

func makeQuestions(opts gameOptions) []question {
	var questions []question

	for _, table := range kana.Table {
		row := table.Basic.Monographs
		hints := getHints(row)
		questions = appendRowQuestions(row, hints, opts, questions)
		// Play only one row in test mode.
		if opts.Test {
			break
		}
	}

	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return questions
}

func getHints(row kana.KanaRow) []string {
	var hints []string
	for _, row := range row {
		if row.Hiragana == "" {
			continue
		}
		var nextHint string
		if row.Alt != "" {
			nextHint = row.Alt
		} else {
			nextHint = row.Romaji
		}
		hints = append(hints, nextHint)
	}
	return hints
}

func appendRowQuestions(
	row kana.KanaRow, hints []string, opts gameOptions, questions []question) []question {

	for _, row := range row {
		if row.Hiragana == "" {
			continue
		}
		var q question
		switch opts.syllabary {
		case hiragana:
			q.kana = row.Hiragana
		case katakana:
			q.kana = row.Katakana
		}
		q.romaji = []string{row.Romaji}
		if row.Alt != "" {
			q.romaji = append(q.romaji, row.Alt)
		}
		q.hints = hints
		questions = append(questions, q)
	}
	return questions
}

func makeGameInput(style lipgloss.Style) textinput.Model {

	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 5
	ti.Width = 5
	ti.Prompt = ""
	ti.TextStyle = style

	return ti
}

func makeProgressBar(opts progressOpts) progress.Model {

	pBar := progress.New(opts.options...)
	if opts.empty != "" {
		pBar.EmptyColor = opts.empty
	}
	if opts.full != "" {
		pBar.FullColor = opts.full
	}
	pBar.PercentageStyle = opts.percStyle
	pBar.ShowPercentage = opts.showPerc

	pBar.Width = opts.width

	return pBar
}

func makeGameKeys() gameKeys {

	return gameKeys{
		Show: key.NewBinding(
			key.WithKeys(";"),
			key.WithHelp(";", "toggle keys"),
		),
		Accept: key.NewBinding(
			key.WithKeys(" ", tea.KeyEnter.String()),
			key.WithHelp("space", "accept"),
		),
		Menu: key.NewBinding(
			key.WithKeys(tea.KeyCtrlO.String()),
			key.WithHelp("ctrl+o", "back to main menu"),
		),
		Hint: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "gimme a hint"),
		),
		Command: key.NewBinding(
			key.WithKeys(":"),
			key.WithHelp(":", "cmd mode"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", tea.KeyEsc.String(), tea.KeyCtrlC.String()),
			key.WithHelp("q", "quit"),
		),
	}
}
