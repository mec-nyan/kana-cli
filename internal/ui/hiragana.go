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

		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "backspace":
			m.input = ""
		case "enter", " ":
			m.tries++
			if m.input == q.romaji[0] {
				m.current++
			}
			m.input = ""
		default:
			m.input += msg.String()
		}
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

	s := fmt.Sprintf("Hiragana: %s\n\n", q.hiragana)

	s += fmt.Sprintf("Write in romaji: %s_\n\n", m.input)

	s += "Hint: (...)\n\n"

	s += fmt.Sprintf("Progress: %.1f%% - Accuracy: %.1f%%\n\n", progress, accuracy)

	return s
}
