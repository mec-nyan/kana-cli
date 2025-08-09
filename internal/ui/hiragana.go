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
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			m.current++
		}
	}
	return m, nil
}

func (m KanaModel) View() string {
	if m.current == len(m.Questions) {
		return "All done!"
	}

	q := m.Questions[m.current]

	progress := 100.0 / float64(len(m.Questions)) * float64(m.current)

	s := fmt.Sprintf("Hiragana: %s\n\n", q.hiragana)

	s += "Write in romaji: _\n\n"

	s += "Hint: (...)\n\n"

	s += fmt.Sprintf("Progress: %.1f%%\n\n", progress)

	return s
}
