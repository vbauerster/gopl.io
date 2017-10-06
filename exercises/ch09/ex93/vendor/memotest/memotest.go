package memotest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func httpGetBody(ctx context.Context, key string) (interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, key, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

var HTTPGetBody = httpGetBody

func incomingURLs() <-chan string {
	ch := make(chan string)
	go func() {
		for _, url := range [...]string{
			"https://golang.org",      // 1
			"https://godoc.org",       // 2
			"https://play.golang.org", // 3
			"http://gopl.io",          // 4
			"https://golang.org",      // 5
			"https://godoc.org",       // 6
			"https://play.golang.org", // 7
			"http://gopl.io",          // 8
		} {
			ch <- url
		}
		close(ch)
	}()
	return ch
}

type M interface {
	Get(ctx context.Context, key string, id int) (interface{}, error)
}

func Concurrent(t *testing.T, m M) {
	var n sync.WaitGroup
	var count int
	for url := range incomingURLs() {
		n.Add(1)
		count++
		go func(url string, count int) {
			defer n.Done()
			var (
				ctx    = context.Background()
				cancel context.CancelFunc
			)
			ctx, cancel = context.WithTimeout(ctx, 2200*time.Millisecond)
			defer cancel()
			start := time.Now()
			value, err := m.Get(ctx, url, count)
			if err != nil {
				fmt.Printf("[%d]:err %v, %s\n", count, err, time.Since(start))
				return
			}
			fmt.Printf("[%d]:ok %s, %s, %d bytes\n", count,
				url, time.Since(start), len(value.([]byte)))
		}(url, count)
	}
	n.Wait()
}
