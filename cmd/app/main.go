package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mec-nyan/kana-cli/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ups! Shit happens dude!\n")
		os.Exit(1)
	}
}
