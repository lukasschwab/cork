package main

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsevents"
)

var (
	configStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	eventStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

func printConfigHeading() {
	pwd, _ := filepath.Abs(".")
	header := fmt.Sprintf("Relative to %s:", pwd)
	println(configStyle.Bold(true).Render(header))
}

func printConfigPair(pair configPair) {
	fmt.Println(configStyle.Render(pair.String()))
}

func printEvent(e fsevents.Event, cmdString string) {
	summary := fmt.Sprintf("[%d] %s → %s", e.ID, e.Path, cmdString)
	fmt.Println(eventStyle.Render(summary))
}
