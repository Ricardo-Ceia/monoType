package main

import (
	"bufio"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type model struct {
	textarea textarea.Model
	text     string
	cursor   int
}

func initialModel() model {
	ta := textarea.New()
	ta.Focus()

	return model{
		textarea: ta,
		text:     randomizeQuotes(getAllWords(readTextFromFile("qoutes.txt")), 30, time.Now().UnixNano()),
		cursor:   0,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) View() string {
	return m.textarea.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
	default:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}
}

func getAllWords(text string) []string {
	return strings.Split(text, " ")
}

func randomizeQuotes(words []string, maxIdx int, seed int64) string {
	r := rand.New(rand.NewSource(seed))
	nums := make([]int, maxIdx)
	for i := range nums {
		nums[i] = i
	}
	r.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})

	var result strings.Builder
	for _, n := range nums {
		result.WriteString(words[n] + " ")
	}
	return strings.TrimSpace(result.String())
}

func readTextFromFile(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	var builder strings.Builder
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		builder.WriteString(line)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panic(err)
		}
	}
	// remove the \n char at the end of the text
	text := builder.String()
	text = strings.TrimSuffix(text, "\n")

	return text
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
