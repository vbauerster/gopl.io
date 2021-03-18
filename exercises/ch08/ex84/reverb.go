// Exercise 8.4: Modify the reverb2 server to use a sync.WaitGroup per
// connection to count the number of active echo goroutines. When it falls to
// zero, close the write half of the TCP connection as described in Exercise
// 8.3. Verify that your modified netcat3 client from that exercise waits for
// the final echoes of multiple concurrent shouts, even after the standard input
// has been closed.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type closeWriter interface {
	CloseWrite() error
}

//!+
func handleConn(c net.Conn) {
	var wg sync.WaitGroup
	count := make(chan int)
	input := bufio.NewScanner(c)
	for input.Scan() {
		wg.Add(1)
		// echo goroutine
		go func(shout string, delay time.Duration) {
			defer wg.Done()
			fmt.Fprintln(c, "\t", strings.ToUpper(shout))
			time.Sleep(delay)
			fmt.Fprintln(c, "\t", shout)
			time.Sleep(delay)
			fmt.Fprintln(c, "\t", strings.ToLower(shout))
			count <- 1
		}(input.Text(), 1*time.Second)
	}
	go func() {
		wg.Wait()
		close(count)
	}()
	var total int
	for i := range count {
		total += i
	}
	fmt.Fprintf(c, "\tTotal: %d\n", total)
	// NOTE: ignoring potential errors from input.Err()
	if cw, ok := c.(closeWriter); ok {
		cw.CloseWrite()
	}
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
