package ui

import (
	"time"

	"github.com/Ricardo-Ceia/monoType/quotes"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		if m.Mode == "menu" {
			return m.handleMenuInput(msg)
		} else if m.Mode == "typping" {
			return m.handleTypingInput(msg)
		} else {
			return m.handleStatsInput(msg)
		}
	case tickMsg:
		if m.Mode == "typping" && !m.StartTime.IsZero() {
			deadline := m.StartTime.Add(time.Duration(m.TimeLimit) * time.Second)
			if time.Now().After(deadline) {
				m.Mode = "stats"
				m.TypedText = ""
				m.Cursor = 0
				m.CorrectChars = 0
				m.StartTime = time.Time{}
				return m, nil
			}
		}
		return m, tickCmd()
	}
	return m, nil
}

func (m Model) handleTypingInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyBackspace:
		if m.Cursor > 0 {
			m.TypedText = m.TypedText[:len(m.TypedText)-1]
			m.Cursor--
		}
	case tea.KeySpace:
		m.TypedText += " "
		m.Cursor++
		if m.Cursor <= len(m.TargetText) {
			if m.TargetText[m.Cursor-1] == ' ' {
				m.CorrectChars++
			}
		}
	case tea.KeyEsc:
		m.Mode = "menu"
		m.SelectedMenu = 0
		m.TypedText = ""
		m.Cursor = 0
		m.CorrectChars = 0
		m.StartTime = time.Time{}
	case tea.KeyRunes:
		if m.StartTime.IsZero() {
			m.StartTime = time.Now()
		}
		for _, r := range msg.Runes {
			m.TypedText += string(r)
			m.Cursor++
			if m.Cursor <= len(m.TargetText) {
				if string(m.TargetText[m.Cursor-1]) == string(r) {
					m.CorrectChars++
				}
			}
		}
	}
	return m, nil
}

func (m Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.SelectedMenu > 0 {
			m.SelectedMenu--
		}
	case tea.KeyDown:
		if m.SelectedMenu < 2 {
			m.SelectedMenu++
		}
	case tea.KeyLeft, tea.KeyRight:
		switch m.SelectedMenu {
		case 0:
			if msg.Type == tea.KeyLeft {
				if m.TimeLimit > 10 {
					m.TimeLimit -= 10
				}
			} else {
				m.TimeLimit += 10
			}
		}
	case tea.KeyEnter:
		if m.SelectedMenu == 1 {
			m.Mode = "typping"
			m.TargetText = quotes.TyppingText(30)
			return m, tickCmd()
		}
	}
	return m, nil
}

func (m Model) handleStatsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.Mode = "menu"
		m.SelectedMenu = 0
		m.TypedText = ""
		m.Cursor = 0
		m.CorrectChars = 0
		m.StartTime = time.Time{}
	case tea.KeyCtrlR:
		m.Mode = "typping"
		m.TargetText = quotes.TyppingText(30)
		return m, tickCmd()
	}
	return m, nil
}
