package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

var program *tea.Program

func main() {
	m := NewFileSelectionModel()
	program = tea.NewProgram(&m, tea.WithAltScreen())
	program.Run()
	fmt.Println("Done")
}
