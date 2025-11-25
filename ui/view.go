package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
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
	mutedColor := lipgloss.Color("#6B7280")

	wpm := 0.0
	if m.TimeLimit > 0 {
		wpm = float64(m.CorrectChars) / 5.0 / (float64(m.TimeLimit) / 60.0)
	}

	accuracy := 0.0
	totalTyped := m.CorrectChars + m.IncorrectChars
	if totalTyped > 0 {
		accuracy = float64(m.CorrectChars) / float64(totalTyped) * 100
	}

	consistency := 0.0
	if len(m.WPMHistory) > 1 {
		var sum float64
		for _, s := range m.WPMHistory {
			sum += s.WPM
		}
		mean := sum / float64(len(m.WPMHistory))
		var sumSquares float64
		for _, s := range m.WPMHistory {
			sumSquares += (s.WPM - mean) * (s.WPM - mean)
		}
		variance := sumSquares / float64(len(m.WPMHistory))
		// Newton's method for square root
		stdDev := 0.0
		if variance > 0 {
			stdDev = variance
			for i := 0; i < 10; i++ {
				stdDev = (stdDev + variance/stdDev) / 2
			}
		}
		if mean > 0 {
			// Coefficient of variation as percentage
			cv := (stdDev / mean) * 100
			consistency = 100 - cv
			if consistency < 0 {
				consistency = 0
			}
		}
	}

	dimStyle := lipgloss.NewStyle().Foreground(mutedColor)

	statsLine := fmt.Sprintf("wpm %.0f  |  acc %.1f%%  |  con %.1f%%", wpm, accuracy, consistency)

	graphWidth := m.Width - 10
	if graphWidth < 30 {
		graphWidth = 30
	}
	if graphWidth > 60 {
		graphWidth = 60
	}
	graphHeight := m.Height - 10
	if graphHeight < 5 {
		graphHeight = 5
	}
	if graphHeight > 12 {
		graphHeight = 12
	}

	var graphContent strings.Builder
	graphContent.WriteString("WPM\n")

	if len(m.WPMHistory) > 0 {
		maxWPM := 0.0
		for _, s := range m.WPMHistory {
			if s.WPM > maxWPM {
				maxWPM = s.WPM
			}
		}
		if maxWPM < 10 {
			maxWPM = 10
		}
		maxWPM = float64(int(maxWPM/10)+1) * 10

		ySteps := 4
		if graphHeight < 8 {
			ySteps = 2
		}

		xSteps := 4
		if graphWidth < 40 {
			xSteps = 2
		}

		pointRows := make([]int, graphWidth)
		maxTime := float64(m.TimeLimit)
		for col := 0; col < graphWidth; col++ {
			// Map column to time value (0 to TimeLimit)
			t := (float64(col) / float64(graphWidth-1)) * maxTime
			
			var wpmVal float64
			if len(m.WPMHistory) == 1 {
				wpmVal = m.WPMHistory[0].WPM
			} else {
				var prev, next WPMSample
				foundPrev := false
				for i := len(m.WPMHistory) - 1; i >= 0; i-- {
					if m.WPMHistory[i].Time <= t {
						prev = m.WPMHistory[i]
						foundPrev = true
						if i+1 < len(m.WPMHistory) {
							next = m.WPMHistory[i+1]
						} else {
							next = prev
						}
						break
					}
				}
				if !foundPrev {
					wpmVal = m.WPMHistory[0].WPM
				} else if prev.Time == next.Time || next.Time == prev.Time {
					wpmVal = prev.WPM
				} else {
					ratio := (t - prev.Time) / (next.Time - prev.Time)
					wpmVal = prev.WPM + ratio*(next.WPM-prev.WPM)
				}
			}
			
			pointRows[col] = int((wpmVal / maxWPM) * float64(graphHeight))
			if pointRows[col] > graphHeight {
				pointRows[col] = graphHeight
			}
			if pointRows[col] < 0 {
				pointRows[col] = 0
			}
		}

		for row := graphHeight; row >= 1; row-- {
			yVal := (float64(row) / float64(graphHeight)) * maxWPM
			showLabel := false
			for step := 0; step <= ySteps; step++ {
				stepRow := (step * graphHeight) / ySteps
				if stepRow == 0 {
					stepRow = 1
				}
				if row == stepRow || row == graphHeight {
					showLabel = true
					break
				}
			}
			if showLabel {
				graphContent.WriteString(dimStyle.Render(fmt.Sprintf("%3.0f│", yVal)))
			} else {
				graphContent.WriteString(dimStyle.Render("   │"))
			}

			for col := 0; col < graphWidth; col++ {
				pointRow := pointRows[col]
				if pointRow == row-1 {
					graphContent.WriteString("•")
				} else if col > 0 {
					prevRow := pointRows[col-1]
					currRow := pointRows[col]
					minR, maxR := prevRow, currRow
					if minR > maxR {
						minR, maxR = maxR, minR
					}
					if row-1 >= minR && row-1 <= maxR {
						graphContent.WriteString("•")
					} else {
						graphContent.WriteString(" ")
					}
				} else {
					graphContent.WriteString(" ")
				}
			}
			graphContent.WriteString("\n")
		}

		graphContent.WriteString(dimStyle.Render("   └" + strings.Repeat("─", graphWidth) + "\n"))

		xAxisLine := strings.Repeat(" ", graphWidth+4)
		xAxisRunes := []rune(xAxisLine)
		for step := 0; step <= xSteps; step++ {
			pos := (step * (graphWidth - 1)) / xSteps
			timeVal := (step * m.TimeLimit) / xSteps
			label := fmt.Sprintf("%d", timeVal)
			startPos := pos + 4
			if step == 0 {
				startPos = 4			}
			for j, ch := range label {
				if startPos+j < len(xAxisRunes) {
					xAxisRunes[startPos+j] = ch
				}
			}
		}
		graphContent.WriteString(dimStyle.Render(string(xAxisRunes) + "s"))
	} else {
		graphContent.WriteString(dimStyle.Render("  No data"))
	}

	footer := dimStyle.Render("Ctrl+R retry  |  Esc menu")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		statsLine,
		"",
		graphContent.String(),
		"",
		footer,
	)

	return content
}
