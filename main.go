package main

import (
	"fmt"
	"github.com/Ricardo-Ceia/monoType/quotes"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
)

type model struct {
	targetText   string
	typedText    string
	cursor       int
	correctChars int
	correctWords int
}

func initialModel() model {
	return model{
		targetText:   quotes.TyppingText(30),
		typedText:    "",
		cursor:       0,
		correctChars: 0,
		correctWords: 0,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyBackspace:
			if m.cursor > 0 {
				m.typedText = m.typedText[:len(m.typedText)-1]
				m.cursor--
			}
		case tea.KeySpace:
			m.typedText += " "
			m.cursor++
			if m.cursor <= len(m.targetText) {
				if m.targetText[m.cursor-1] == ' ' {
					m.correctChars++
				}
			}
		case tea.KeyRunes:
			for _, r := range msg.Runes {
				m.typedText += string(r)
				m.cursor++
				if m.cursor <= len(m.targetText) {
					if string(m.targetText[m.cursor-1]) == string(r) {
						m.correctChars++
					}
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var display strings.Builder
	for i, ch := range m.targetText {
		if i < len(m.typedText) {
			if rune(m.targetText[i]) == rune(m.typedText[i]) {
				display.WriteString(fmt.Sprintf("\033[32m%c\033[0m", ch)) // Green for correct
			} else {
				display.WriteString(fmt.Sprintf("\033[31m%c\033[0m", ch)) // Red for wrong
			}
		} else if i == len(m.typedText) {
			display.WriteString(fmt.Sprintf("\033[1m\033[4m%c\033[0m", ch)) // Bold + underline cursor
		} else {
			display.WriteRune(ch)
		}
	}
	stats := fmt.Sprintf("\n\nTyped: %d/%d | Correct: %d",
		len(m.typedText), len(m.targetText), m.correctChars)
	return display.String() + stats
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
