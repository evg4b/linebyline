# Line by line writer [![Go Report Card](https://goreportcard.com/badge/github.com/evg4b/linebyline)](https://goreportcard.com/report/github.com/evg4b/linebyline)

Thread-safe `io.WriterCloser` allowing you to write from different sources in one line by line. 
Writing to the source occurs only after the `end of the line` or the `closing of the generated io.WriterCloser` or `closing the WriterGroup`.

# Example
``` GO
package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/evg4b/linebyline"
)

func main() {
	group := linebyline.NewWriterGroup(os.Stdout)

	wg := sync.WaitGroup{}

	wr1 := group.CreateWriter()
	wr2 := group.CreateWriter()

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			fmt.Fprintf(wr1, "[#1] line %d ", i)
			// do something else ...
			time.Sleep(10 * time.Millisecond)

			fmt.Fprintln(wr1, "first writer")
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			fmt.Fprintf(wr2, "[#2] line %d ", i)
			// do something else ...
			time.Sleep(20 * time.Millisecond)

			fmt.Fprintln(wr2, "second writer")
		}
	}()

	wg.Wait()

	err := group.Close()
	if err != nil {
		panic(err)
	}
}
```
Output: 
```
[#1] line 0 first writer
[#1] line 1 first writer 
[#2] line 0 second writer
[#1] line 2 first writer
[#1] line 3 first writer 
[#2] line 1 second writer
[#1] line 4 first writer
[#2] line 2 second writer
[#1] line 12 first writer
[#1] line 13 first writer
[#2] line 6 second writer
[#2] line 7 second writer
...
```