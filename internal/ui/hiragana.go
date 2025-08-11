package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mec-nyan/kana-master/pkg/kana"
)

/*
This UI will look something like this:

Hiragana: あ

Write in romaji: _

Hint: (a i u e o)

Progress: x% ||||||||||||||||____

*/

const (
	padding = 4
)

type (
	question struct {
		hiragana string
		romaji   []string
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
	}

	errMsg error
)

func KanaInitialModel(testMode bool) tea.Model {
	var questions []question
	// TODO: Shuffle
	for _, table := range kana.Table {
		// For now, only monographs
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
			questions = append(questions, q)
		}

		// Play only one row in test mode.
		if testMode {
			break
		}
	}

	ti := textinput.New()
	ti.Placeholder = "..."
	ti.Focus()
	ti.CharLimit = 5
	ti.Width = 5
	ti.Prompt = ""

	prog := progress.New()

	return KanaModel{
		Questions: questions,
		textInput: ti,
		progress:  prog,
	}
}

func (m KanaModel) Init() tea.Cmd {
	return nil
}

func (m KanaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	q := m.Questions[m.current]

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2
		return m, nil

	case tea.KeyMsg:

		switch msg.String() {
		// There's no hiragana/katakana that starts with "q".
		// We can safely use this to quit.
		case tea.KeyCtrlC.String(), tea.KeyEsc.String(), tea.KeyCtrlD.String(), "q":
			m.quit = true
			return m, tea.Quit

		// Use both enter or space to accept input value.
		case tea.KeyEnter.String(), tea.KeySpace.String():
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
	paddingLeft := strings.Repeat(" ", padding)

	progress := 100.0 / float64(len(m.Questions)) * float64(m.current)
	accuracy := 0.0
	if m.tries > 0 {
		accuracy = float64(m.current) / float64(m.tries) * 100.0
	}

	// TODO: Use lipgloss for padding and styling!
	s := fmt.Sprintf(`
%s%s

%sHiragana: %s

%sWrite in romaji: %s

%sHint: (...)

%sProgress: %0.1f%% - Accuracy: %0.1f%%

%s[3;38:5:8mPress <esc> to quit.[0m`,
		paddingLeft, m.progress.ViewAs(progress/100),
		paddingLeft, q.hiragana,
		paddingLeft, m.textInput.View(),
		paddingLeft,
		paddingLeft, progress, accuracy,
		paddingLeft)

	return s
}
