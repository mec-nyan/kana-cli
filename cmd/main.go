package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mec-nyan/kana-cli/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	if model, err := p.Run(); err != nil {
		fmt.Printf("Ups! Shit happens dude!\n")
		os.Exit(1)
	} else {
		fmt.Printf("You've selected:\n\n")
		m := model.(ui.Model)
		for _, step := range m.Steps {
			fmt.Printf("  %s -> ", step.Name)
			for _, choice := range step.Options {
				if choice.Selected {
					fmt.Printf("%s\n", choice.Name)
				}
			}
		}
	}
}
