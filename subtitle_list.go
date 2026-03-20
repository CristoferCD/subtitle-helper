package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type subtitleListModel struct {
	selectedFile       string
	streams            list.Model
	selectedItem       FfmpegStream
	chosen             bool
	extracted          bool
	extractionProgress progress.Model
	windowWidth        int
	windowHeight       int
}

type listStreamsMsg struct {
	streams []FfmpegStream
}

func (s FfmpegStream) Title() string {
	return fmt.Sprintf("Subtitle #%d", s.Index)
}

func (s FfmpegStream) Description() string {
	return fmt.Sprintf(" %s | Lang: %s", s.CodecLongName, s.Tags["language"])
}

func (s FfmpegStream) FilterValue() string {
	return fmt.Sprintf("%d%s%s", s.Index, s.CodecLongName, s.Tags["language"])
}

type extractedSubtitleMsg struct {
	srtFile string
}

type extractionProgressMsg struct {
	progress float64
}

func (m subtitleListModel) Init() tea.Cmd {
	return m.listSubtitleStreams(m.selectedFile)
}

func (m subtitleListModel) listSubtitleStreams(file string) tea.Cmd {
	return func() tea.Msg {
		streams := ListSubtitleStreams(file)
		return listStreamsMsg{streams}
	}
}

func (m subtitleListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case listStreamsMsg:
		items := make([]list.Item, len(msg.streams))
		for idx, sub := range msg.streams {
			items[idx] = sub
		}
		m.streams = list.New(items, list.NewDefaultDelegate(), m.windowWidth, m.windowHeight)
		m.streams.Title = "Available streams"
	case extractedSubtitleMsg:
		m.chosen = false
		m.extracted = true
		return m, translate(msg.srtFile)
	case extractionProgressMsg:
		return m, m.extractionProgress.SetPercent(msg.progress / 100)
	case progress.FrameMsg:
		progressModel, cmd := m.extractionProgress.Update(msg)
		m.extractionProgress = progressModel.(progress.Model)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			prevModel := NewFileSelectionModel()
			return prevModel, prevModel.Init()
		case "enter":
			selected, ok := m.streams.SelectedItem().(FfmpegStream)
			if ok {
				m.chosen = true
				m.selectedItem = selected
				return m, extract(m.selectedFile, m.selectedItem)
			}
		}
	case tea.WindowSizeMsg:
		m.extractionProgress.Width = msg.Width - 8
		m.streams.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.streams, cmd = m.streams.Update(msg)
	return m, cmd
}

func (m subtitleListModel) View() string {
	if m.chosen {
		return "Extracting... " + m.extractionProgress.View()
	} else if m.extracted {
		return "Translating..."
	} else {
		return m.streams.View()
	}
}

func extract(file string, sub FfmpegStream) tea.Cmd {
	return func() tea.Msg {
		progressWriter := func(progress float64) {
			program.Send(extractionProgressMsg{progress})
		}
		file := ExtractSubtitle(file, sub, progressWriter)
		return extractedSubtitleMsg{file}
	}
}

func translate(subtitle string) tea.Cmd {
	return func() tea.Msg {
		Translate(subtitle)
		return tea.Quit()
	}
}
