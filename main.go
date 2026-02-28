package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	fp := filepicker.New()
	fp.AllowedTypes = []string{".mkv"}
	fp.CurrentDirectory = "."
	fp.DirAllowed = true

	m := model{
		filepicker: fp,
	}
	tea.NewProgram(&m, tea.WithAltScreen()).Run()
	fmt.Println("Done")
}
