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

	// Auto mode is going to an Option in the menu (to be implemented).
	// In Auto mode, there's no need to press enter or space to accept the input.
	autoModeLong := flag.Bool("auto", false, "Play in auto mode (options).")
	autoModeShort := flag.Bool("a", false, "Play in auto mode (options) (short).")

	flag.Parse()

	testMode := *testModeShort || *testModeLong
	autoMode := *autoModeShort || *autoModeLong

	opts := ui.CLIOptions{
		Test: testMode,
		Auto: autoMode,
	}

	p := tea.NewProgram(ui.MainMenuInitialModel(opts))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ups! Shit happens dude!\n")
		os.Exit(1)
	}
}
