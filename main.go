package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var program *tea.Program

var config Config

type Config struct {
	workDir            string
	libreTranslateHost string
	libreTranslatePort string
}

func LoadConfig() Config {
	workDir := os.Getenv("WORK_DIRECTORY")
	if len(workDir) == 0 {
		workDir = "."
	}

	translateHost := os.Getenv("LIBRETRANSLATE_HOST")
	if len(workDir) == 0 {
		translateHost = "libretranslate"
	}

	translatePort := os.Getenv("LIBRETRANSLATE_PORT")
	if len(workDir) == 0 {
		translatePort = "5000"
	}

	return Config{
		workDir:            workDir,
		libreTranslateHost: translateHost,
		libreTranslatePort: translatePort,
	}
}

func main() {
	config = LoadConfig()

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
