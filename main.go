package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}


func cleanInput(text string) []string {
	cleanedText := strings.TrimSpace(strings.ToLower(text))
	words := strings.Fields(cleanedText)

	return words
}