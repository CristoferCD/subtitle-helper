package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var program *tea.Program

func main() {
	m := NewFileSelectionModel()
	program = tea.NewProgram(&m, tea.WithAltScreen())

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := program.Run(); err != nil {
		log.Printf("Alas, there's been an error: %v", err.Error())
		os.Exit(1)
	}
}
