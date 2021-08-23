# Line by line writer [![Go Report Card](https://goreportcard.com/badge/github.com/evg4b/linebyline)](https://goreportcard.com/report/github.com/evg4b/linebyline)

Thread-safe `io.WriterCloser` allowing you to write from different sources in one line by line. 
Writing to the source occurs only after the `end of the line` or the `closing of the generated io.WriterCloser` or `closing the WriterGroup`.

# Example
``` GO
package main

import (
	"fmt"
	"linebyline"
	"os"
	"sync"
	"time"
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
			fmt.Fprintf(wr1, "[#1] first writer line %d\n", i)
            // do something else ...
			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			fmt.Fprintf(wr2, "[#2] second writer line %d\n", i)
            // do something else ...
			time.Sleep(20 * time.Millisecond)
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
[#2] second writer line 0
[#1] first writer line 0
[#1] first writer line 1 
[#2] second writer line 1
[#1] first writer line 2
[#1] first writer line 3
[#2] second writer line 2
[#1] first writer line 4
[#1] first writer line 5 
[#2] second writer line 3
[#1] first writer line 6
[#1] first writer line 7 
[#2] second writer line 4
...
```