package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mec-nyan/kana-cli/internal/ui/hiragana"
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

	// TODO: Actually, there's gonna be an options menu that will handle this.
	p := tea.NewProgram(hiragana.KanaInitialModel(hiragana.Options{
		Test: testMode, Auto: autoMode}))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ups! Shit happens dude!\n")
		os.Exit(1)
	}
}
