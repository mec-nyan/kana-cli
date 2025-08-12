package ui

import (
	"fmt"
	"strings"
	"time"

	. "github.com/mec-nyan/kana-cli/internal/palette"
	"github.com/mec-nyan/kana-master/pkg/kana"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxBarWidth = 80
)

var (
	highlightStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(Mauve))
	inputStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(Teal))
	hintStyle      = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(Overlay0))
)

type (
	Options struct {
		Test bool
		Auto bool
	}

	Question struct {
		hiragana string
		romaji   []string
		hints    []string
		played   bool
	}

	KanaModel struct {
		Questions []Question
		current   int
		textInput textinput.Model
		tries     int
		percent   int
		Progress  progress.Model
		quit      bool
		err       error
		// TODO: Not implemented yet!
		// Add a menu entry to select autoMode "on/off".
		// In autoMode, you don't need to press enter or space,
		// your input is compared with the current kana each time and move
		// to the next question as soon as it it correct.
		autoMode bool
		keys     kanaKeyMap
		help     help.Model
		style    lipgloss.Style
		hint     bool
		end      bool
		accuracy float64
		menu     MainMenuModel
	}

	kanaKeyMap struct {
		Show    key.Binding
		Accept  key.Binding
		Menu    key.Binding
		Command key.Binding
		Hint    key.Binding
		Help    key.Binding
		Quit    key.Binding
	}

	errMsg error

	tickMsg struct{}
)

func (k kanaKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Show}
}

func (k kanaKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Show, k.Accept, k.Menu, k.Hint, k.Command, k.Help, k.Quit},
	}
}

func KanaInitialModel(menu MainMenuModel, opts Options) tea.Model {
	var questions []Question
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
			var q Question
			q.hiragana = row.Hiragana
			q.romaji = []string{row.Romaji}
			if row.Alt != "" {
				q.romaji = append(q.romaji, row.Alt)
			}
			q.hints = hints
			questions = append(questions, q)
		}

		// Play only one row in test mode.
		if opts.Test {
			break
		}
	}

	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 5
	ti.Width = 5
	ti.Prompt = ""
	ti.TextStyle = inputStyle

	prog := progress.New(progress.WithGradient(Mauve, Sapphire), progress.WithFillCharacters('▂', '▂'))
	prog.EmptyColor = Surface0
	prog.PercentageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(Subtext0))

	keys := kanaKeyMap{
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
		Progress:  prog,
		keys:      keys,
		help:      help.New(),
		style:     appStyle,
		autoMode:  opts.Auto,
		menu:      menu,
	}
}

func (m KanaModel) Init() tea.Cmd {
	return nil
}

func (m KanaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var q Question
	if m.current < len(m.Questions) {
		q = m.Questions[m.current]
	}

	switch msg := msg.(type) {

	case tickMsg:
		m.quit = true
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.style = m.style.Width(msg.Width)
		m.Progress.Width = min(msg.Width-padding*2, maxBarWidth)
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
			float64(len(m.Questions))/float64(m.tries)*100,
			m.help.View(m.keys))

		return m.style.Render(s)
	}

	if m.current >= len(m.Questions) {
		return "Nope!"
	}

	q := m.Questions[m.current]

	hint := "Hint: ..."
	if m.hint {
		hint = "Hint: " + strings.Join(q.hints, ", ")
	}

	progress := 100.0 / float64(len(m.Questions)) * float64(m.current)
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
		highlightStyle.Render(q.hiragana),
		m.textInput.View(),
		hintStyle.Render(hint),
		m.accuracy,
		m.help.View(m.keys))

	return m.style.Render(s)
}

// TODO: How to check in autoMode?
func (m KanaModel) checkAnswer(q Question) (tea.Model, tea.Cmd) {
	m.hint = false
	m.tries++
	guess := m.textInput.Value()
	for _, rmj := range q.romaji {
		if guess == rmj {
			m.current++
			if m.current == len(m.Questions) {
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
