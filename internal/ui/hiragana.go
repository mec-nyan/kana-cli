package ui

import (
	"fmt"

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

type question struct {
	hiragana string
	romaji   []string
	played   bool
}

type KanaModel struct {
	Questions []question
	current   int
	input     string
	tries     int
	quit      bool
}

func KanaInitialModel() tea.Model {
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
	}
	return KanaModel{
		Questions: questions,
	}
}

func (m KanaModel) Init() tea.Cmd {
	return nil
}

func (m KanaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	q := m.Questions[m.current]

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch s := msg.String(); s {
		case "ctrl+c", "esc":
			m.quit = true
			return m, tea.Quit
		case "backspace":
			m.input = ""
		case "enter", " ":
			m.tries++
			guess := m.input
			// There may be different romaji associated to this kana.
			for _, rmj := range q.romaji {
				if guess == rmj {
					m.current++
				}
			}
			m.input = ""
		default:
			if len(s) == 1 && s[0] >= 'a' && s[0] <= 'z' {
				m.input += s
			}
		}
	}

	if m.current == len(m.Questions) {
		return m, tea.Quit
	}

	return m, nil
}

func (m KanaModel) View() string {
	if m.current == len(m.Questions) {
		return "All done!"
	}

	if m.quit {
		return "Good bye then!"
	}

	q := m.Questions[m.current]

	progress := 100.0 / float64(len(m.Questions)) * float64(m.current)
	accuracy := 0.0
	if m.tries > 0 {
		accuracy = float64(m.current) / float64(m.tries) * 100.0
	}

	s := fmt.Sprintf("\n\tHiragana: %s\n\n", q.hiragana)

	s += fmt.Sprintf("\tWrite in romaji: %s_\n\n", m.input)

	s += "\tHint: (...)\n\n"

	s += fmt.Sprintf("\tProgress: %.1f%% - Accuracy: %.1f%%\n\n", progress, accuracy)

	s += "\t[3;38:5:8mPress <esc> to quit.[0m"
	return s
}
