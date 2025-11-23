package main

import (
	"fmt"
	"github.com/Ricardo-Ceia/monoType/quotes"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
	"time"
)

type model struct {
	mode         string
	targetText   string
	typedText    string
	selectedMenu int
	timeLimit    int
	startTime    time.Time
	cursor       int
	correctChars int
	correctWords int
}

func initialModel() model {
	return model{
		mode:         "typping",
		targetText:   quotes.TyppingText(30),
		typedText:    "",
		selectedMenu: 0,
		timeLimit:    30,
		startTime:    time.Time{},
		cursor:       0,
		correctChars: 0,
		correctWords: 0,
	}
}

func (m model) Init() tea.Cmd {
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == "menu" {
			return m.handleMenuInput(msg)
		} else {
			return m.handleTypingInput(msg)
		}
	case tickMsg:
		if m.mode == "typping" && !m.startTime.IsZero() {
			deadline := m.startTime.Add(time.Duration(m.timeLimit) * time.Second)
			if time.Now().After(deadline) {
				m.startTime = time.Time{}
				return m, nil
			}
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m model) handleTypingInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	case tea.KeyEsc:
		m.mode = "menu"
		m.selectedMenu = 0
		m.typedText = ""
		m.cursor = 0
		m.correctChars = 0
		m.startTime = time.Time{}
	case tea.KeyRunes:
		if m.startTime.IsZero() {
			m.startTime = time.Now()
		}
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
	return m, nil
}

func (m model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.selectedMenu > 0 {
			m.selectedMenu--
		}
	case tea.KeyDown:
		if m.selectedMenu < 2 {
			m.selectedMenu++
		}
	case tea.KeyLeft, tea.KeyRight:
		switch m.selectedMenu {
		case 0:
			if msg.Type == tea.KeyLeft {
				if m.timeLimit > 10 {
					m.timeLimit -= 10
				}
			} else {
				m.timeLimit += 10
			}
		}
	case tea.KeyEnter:
		if m.selectedMenu == 1 {
			m.mode = "typping"
			m.targetText = quotes.TyppingText(30)
		}
	}
	return m, nil
}

func (m model) View() string {

	if m.mode == "menu" {
		return m.viewMenu()
	} else {
		return m.viewTypper()
	}

}

func (m model) viewTypper() string {
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
	stats := fmt.Sprintf("\n\nTimer:%s Typed: %d/%d | Correct: %d",
		m.viewTimer(), len(m.typedText), len(m.targetText), m.correctChars)
	return display.String() + stats
}

func (m model) viewMenu() string {
	var display strings.Builder
	display.WriteString("Settings and Stats:\n\n")

	if m.selectedMenu == 0 {
		display.WriteString(fmt.Sprintf("> Time limit: %d seconds (← →)\n", m.timeLimit))
	} else {
		display.WriteString(fmt.Sprintf("  Time limit: %d seconds (← →)\n", m.timeLimit))
	}
	if m.selectedMenu == 1 {
		display.WriteString("> Exit Menu\n")
	} else {
		display.WriteString("  Exit Menu\n")
	}

	display.WriteString("\n(↑/↓ to navigate, ← → to change settings, Enter to select, Ctrl+C to quit)\n")
	return display.String()
}

func (m model) viewTimer() string {
	if m.startTime.IsZero() {
		return fmt.Sprintf("TIME: %02d:%02d", m.timeLimit/60, m.timeLimit%60)
	}

	deadline := m.startTime.Add(time.Duration(m.timeLimit) * time.Second)
	remaining := time.Until(deadline)

	if remaining < 0 {
		remaining = 0
	}
	secondsLeft := int(remaining.Seconds())
	return fmt.Sprintf("TIME: %02d:%02d", secondsLeft/60, secondsLeft%60)
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
