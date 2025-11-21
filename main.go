package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

/*
func randonmizeQoutes(text string, seed int) {

}
*/

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
	//remove the \n char at the end of the text
	text := builder.String()
	text = strings.TrimSuffix(text, "\n")

	return text
}

func main() {
	filepath := "qoutes.txt"
	text := readTextFromFile(filepath)
	fmt.Println(text)
}
