package ui

import (
	"fmt"
	"strings"
	"time"
)

func (m Model) View() string {
	if m.Mode == "menu" {
		return m.viewMenu()
	} else if m.Mode == "typping" {
		return m.viewTypper()
	} else {
		return m.viewStats()
	}
}

func (m Model) viewTypper() string {
	var display strings.Builder

	viewportWidth := m.Width - 10

	if viewportWidth < 40 {
		viewportWidth = 40
	}

	viewportStart := m.Cursor - viewportWidth/2

	if viewportStart < 0 {
		viewportStart = 0
	}

	viewportEnd := viewportStart + viewportWidth

	if viewportEnd > len(m.TargetText) {
		viewportEnd = len(m.TargetText)
		if viewportEnd-viewportStart < viewportWidth {
			viewportStart = viewportEnd - viewportWidth
			if viewportStart < 0 {
				viewportStart = 0
			}
		}
	}

	for i := viewportStart; i < viewportEnd; i++ {
		ch := rune(m.TargetText[i])
		if i < len(m.TypedText) {
			if rune(m.TargetText[i]) == rune(m.TypedText[i]) {
				display.WriteString(fmt.Sprintf("\033[32m%c\033[0m", ch)) // Green for correct
			} else {
				display.WriteString(fmt.Sprintf("\033[31m%c\033[0m", ch)) // Red for wrong
			}
		} else if i == len(m.TypedText) {
			display.WriteString(fmt.Sprintf("\033[1m\033[4m%c\033[0m", ch)) // Bold + underline cursor
		} else {
			display.WriteRune(ch)
		}
	}

	stats := fmt.Sprintf("\n\nTimer:%s Typed: %d/%d | Correct: %d",
		m.viewTimer(), len(m.TypedText), len(m.TargetText), m.CorrectChars)
	return display.String() + stats
}

func (m Model) viewMenu() string {
	var display strings.Builder
	display.WriteString("Settings and Stats:\n\n")

	if m.SelectedMenu == 0 {
		display.WriteString(fmt.Sprintf("> Time limit: %d seconds (← →)\n", m.TimeLimit))
	} else {
		display.WriteString(fmt.Sprintf("  Time limit: %d seconds (← →)\n", m.TimeLimit))
	}
	if m.SelectedMenu == 1 {
		display.WriteString("> Exit Menu\n")
	} else {
		display.WriteString("  Exit Menu\n")
	}

	display.WriteString("\n(↑/↓ to navigate, ← → to change settings, Enter to select, Ctrl+C to quit)\n")
	return display.String()
}

func (m Model) viewTimer() string {
	if m.StartTime.IsZero() {
		return fmt.Sprintf("TIME: %02d:%02d", m.TimeLimit/60, m.TimeLimit%60)
	}

	deadline := m.StartTime.Add(time.Duration(m.TimeLimit) * time.Second)
	remaining := time.Until(deadline)

	if remaining < 0 {
		remaining = 0
	}
	secondsLeft := int(remaining.Seconds())
	return fmt.Sprintf("TIME: %02d:%02d", secondsLeft/60, secondsLeft%60)
}

func (m Model) viewStats() string {
	return "-----stats page-----"
}
