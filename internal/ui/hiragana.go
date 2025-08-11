package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mec-nyan/kana-master/pkg/kana"
)

const (
	padding     = 4
	maxBarWidth = 80
)

type (
	question struct {
		hiragana string
		romaji   []string
		hints    []string
		played   bool
	}

	KanaModel struct {
		Questions []question
		current   int
		textInput textinput.Model
		tries     int
		percent   int
		progress  progress.Model
		quit      bool
		err       error
		// TODO: Not implemented yet!
		// Add a menu entry to select autoMode "on/off".
		// In autoMode, you don't need to press enter or space,
		// your input is compared with the current kana each time and move
		// to the next question as soon as it it correct.
		autoMode bool
		keys     keyMap
		help     help.Model
		style    lipgloss.Style
		hint     bool
	}

	keyMap struct {
		Show    key.Binding
		Accept  key.Binding
		Back    key.Binding
		Command key.Binding
		Hint    key.Binding
		Help    key.Binding
		Quit    key.Binding
	}

	errMsg error
)

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Show}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Show, k.Accept, k.Back, k.Hint, k.Command, k.Help, k.Quit},
	}
}

func KanaInitialModel(testMode bool) tea.Model {
	var questions []question
	// TODO: Shuffle
	for _, table := range kana.Table {
		// For now, only monographs

		var hints []string
		for _, row := range table.Basic.Monographs {
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

		for _, row := range table.Basic.Monographs {
			if row.Hiragana == "" {
				continue
			}
			var q question
			q.hiragana = row.Hiragana
			q.romaji = []string{row.Romaji}
			if row.Alt != "" {
				q.romaji = append(q.romaji, row.Alt)
			}
			q.hints = hints
			questions = append(questions, q)
		}

		// Play only one row in test mode.
		if testMode {
			break
		}
	}

	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 5
	ti.Width = 5
	ti.Prompt = ""

	prog := progress.New(progress.WithGradient(mauve, sapphire), progress.WithFillCharacters('▂', '▂'))
	prog.EmptyColor = surface0
	prog.PercentageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(subtext0))

	keys := keyMap{
		Show: key.NewBinding(
			key.WithKeys(";"),
			key.WithHelp(";", "toggle keys"),
		),
		Accept: key.NewBinding(
			key.WithKeys(" ", tea.KeyEnter.String()),
			key.WithHelp("space", "accept"),
		),
		Back: key.NewBinding(
			key.WithKeys(tea.KeyCtrlO.String()),
			key.WithHelp("ctrl+o", "back"),
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

	return KanaModel{
		Questions: questions,
		textInput: ti,
		progress:  prog,
		keys:      keys,
		help:      help.New(),
		style:     lipgloss.NewStyle().Padding(1, padding).Foreground(lipgloss.Color(subtext1)),
	}
}

func (m KanaModel) Init() tea.Cmd {
	return nil
}

func (m KanaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	q := m.Questions[m.current]

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.style = m.style.Width(msg.Width)
		m.progress.Width = min(msg.Width-padding*2, maxBarWidth)
		return m, nil

	case tea.KeyMsg:

		switch {
		// There's no hiragana/katakana that starts with "q".
		// We can safely use this to quit.
		case key.Matches(msg, m.keys.Quit):
			m.quit = true
			return m, tea.Quit

		// Use both enter or space to accept input value.
		case key.Matches(msg, m.keys.Accept):
			m.hint = false
			m.tries++
			guess := m.textInput.Value()
			for _, rmj := range q.romaji {
				if guess == rmj {
					m.current++
					if m.current == len(m.Questions) {
						return m, tea.Quit
					}
				}
			}
			m.textInput.Reset()
			return m, nil

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
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m KanaModel) View() string {
	if m.current == len(m.Questions) {
		return "All done!"
	}

	if m.quit {
		return "Good bye then!"
	}

	q := m.Questions[m.current]

	hint := "Hint: ..."
	if m.hint {
		hint = "Hint: " + strings.Join(q.hints, ", ")
	}

	progress := 100.0 / float64(len(m.Questions)) * float64(m.current)
	accuracy := 0.0
	if m.tries > 0 {
		accuracy = float64(m.current) / float64(m.tries) * 100.0
	}

	// TODO: Use lipgloss for padding and styling!
	s := fmt.Sprintf(`%s

Hiragana: %s

Write in romaji: %s

%s 

Accuracy: %0.1f%%

%s`,
		m.progress.ViewAs(progress/100),
		q.hiragana,
		m.textInput.View(),
		lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(overlay0)).Render(hint),
		accuracy,
		m.help.View(m.keys))

	return m.style.Render(s)
}
