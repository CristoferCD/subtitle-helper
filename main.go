package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewFileSelectionModel()
	tea.NewProgram(&m, tea.WithAltScreen()).Run()
	fmt.Println("Done")
}
