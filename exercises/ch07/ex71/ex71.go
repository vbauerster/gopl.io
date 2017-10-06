// Exercise 7.1: Using the ideas from ByteCounter, implement counters for words
// and for lines. You will find bufio.ScanWords useful.
package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	// Set the split function for the scanning operation
	scanner.Split(bufio.ScanWords)
	// Count the words
	for scanner.Scan() {
		*c++
	}
	return len(p), scanner.Err()
}

func main() {
	// An artificial input source.
	const input = "Now is the winter of our discontent,\nMade glorious summer by this sun of York.\n"

	var c WordCounter
	fmt.Fprintln(&c, input)
	fmt.Println(c) // 15
}
