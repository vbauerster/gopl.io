package main

// Exercise 5.2: Write a function to populate a mapping from element names—p, div,
// span, and so on—to the number of elements with that name in an HTML document
// tree.

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	counts := make(map[string]int)
	mapElements(counts, doc)
	printMap(counts)
}

func mapElements(counts map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		counts[n.Data]++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		mapElements(counts, c)
	}
}

func printMap(counts map[string]int) {
	for k, v := range counts {
		if v > 1 {
			fmt.Printf("%q:\t%.2d\n", k, v)
		}
	}
}
