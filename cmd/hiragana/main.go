package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mec-nyan/kana-cli/internal/ui"
)

func main() {
	testModeLong := flag.Bool("test", false, "Play in test mode.")
	testModeShort := flag.Bool("t", false, "Play in test mode (short).")

	flag.Parse()

	testMode := *testModeShort || *testModeLong

	p := tea.NewProgram(ui.KanaInitialModel(testMode))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ups! Shit happens dude!\n")
		os.Exit(1)
	} }
