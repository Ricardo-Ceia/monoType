package ui

import (
	"time"

	"github.com/Ricardo-Ceia/monoType/quotes"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type WPMSample struct {
	Time float64
	WPM  float64
}

type Model struct {
	Mode           string
	TargetText     string
	TypedText      string
	SelectedMenu   int
	TimeLimit      int
	StartTime      time.Time
	Cursor         int
	CorrectChars   int
	IncorrectChars int
	CorrectWords   int
	TotalChars     int
	Width          int
	Height         int
	WPMHistory     []WPMSample
}

func InitialModel() Model {
	return Model{
		Mode:           "typping",
		TargetText:     quotes.TyppingText(30),
		TypedText:      "",
		SelectedMenu:   0,
		TimeLimit:      30,
		StartTime:      time.Time{},
		Cursor:         0,
		CorrectChars:   0,
		IncorrectChars: 0,
		CorrectWords:   0,
		TotalChars:     0,
		Width:          80,
		Height:         24,
		WPMHistory:     []WPMSample{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		tickCmd(),
	)
}

type tickMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
