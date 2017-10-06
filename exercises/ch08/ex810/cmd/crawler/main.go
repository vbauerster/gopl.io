// Exercise 8.10:
// HTTP requests may be cancelled by closing the optional Cancel channel in the
// http.Request struct. Modify the web crawler of Section 8.6 to support
// cancellation.
// Hint: the http.Get convenience function does not give you an opportunity to
// customize a Request. Instead, create the request using http.NewRequest, set
// its Cancel field, then perform the request by calling
// http.DefaultClient.Do(req).
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"links"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Depth int `long:"depth" default:"2" description:"depth limit"`
}

var validURL = regexp.MustCompile(`^https?://.+`)

// tokens is a counting semaphore used to
// enforce a limit of 8 concurrent requests.
var tokens = make(chan struct{}, 8)

type crawlURL struct {
	url   string
	depth int
}

func (c crawlURL) String() string {
	return c.url
}

func crawl(ctx context.Context, page crawlURL) []crawlURL {
	fmt.Println(page)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(ctx, page.url)
	<-tokens // release the token
	if err != nil {
		log.Print(err)
	}
	children := make([]crawlURL, len(list))
	for i, url := range list {
		children[i].url = url
		children[i].depth = page.depth + 1
	}
	return children
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS] [url1 url2 ...]"
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if len(args) == 0 {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	var urls = make([]crawlURL, len(args))
	for i, arg := range args {
		if !validURL.MatchString(arg) {
			fmt.Fprintf(os.Stderr, "Unsupported url: %s\n", arg)
			os.Exit(1)
		}
		urls[i].url = arg
	}

	if f, err := os.Create("err.log"); err == nil {
		log.SetOutput(f)
		defer f.Close()
	}

	// Create context for further cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println("Press return key to abort!")
	go func() {
		os.Stdin.Read(make([]byte, 1))
		cancel()
	}()

	worklist := make(chan []crawlURL)
	var n int // number of pending sends to worklist

	// Start with the command line arguments
	n++
	go func() { worklist <- urls }()

	// Crawl the web concurrently
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, page := range list {
			if !seen[page.url] && page.depth < opts.Depth {
				seen[page.url] = true
				n++
				go func(page crawlURL) {
					worklist <- crawl(ctx, page)
				}(page)
			}
		}
	}
}
