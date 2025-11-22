package quotes

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

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

func TyppingText(maxIdx int) string {
	textData := readTextFromFile("quotes.txt")
	wordsFromText := getAllWords(textData)
	randomizedQuotes := randomizeQuotes(wordsFromText, maxIdx, time.Now().UnixNano())
	return randomizedQuotes
}
