package main

// Exercise 4.9: Write a program wordfreq to report the frequency of each word
// in an input text file. Call input.Split(bufio.ScanWords) before the first
// call to Scan to break the input into words instead of lines.

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	words := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		words[strings.Trim(input.Text(), ".,!?")]++
	}
	for word, count := range words {
		fmt.Printf("%s \t %d\n", word, count)
	}
}
