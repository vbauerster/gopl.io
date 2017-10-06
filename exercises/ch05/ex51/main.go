package main

// Exercise 5.1: Change the findlinks program to traverse the n.FirstChild linked
// list using recursive calls to visit instead of a loop.

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatalf("findlinks: %v\n", err)
	}
	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

// In-Place Slice Technique (p. 91)
func visit(links []string, n *html.Node) []string {
	if n == nil {
		return links
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	links = visit(links, n.FirstChild)
	links = visit(links, n.NextSibling)
	return links
}
