# Line by line writer [![Go Report Card](https://goreportcard.com/badge/github.com/evg4b/linebyline)](https://goreportcard.com/report/github.com/evg4b/linebyline)

Thread-safe `io.WriterCloser` allowing you to write from different sources in one line by line.
Writing to the source occurs only after the `end of the line` or the `closing of the generated io.WriterCloser`.

## Example:

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

	wg := sync.WaitGroup{}

	wr1 := linebyline.NewByLineWriter(
	    linebyline.WithOutWriter(os.Stdout),
	)
	wr2 := linebyline.NewByLineWriter(
	    linebyline.WithOutWriter(os.Stdout),
	)

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

	err := w1.Close()
	if err != nil {
		panic(err)
	}

	err := w2.Close()
	if err != nil {
		panic(err)
	}
}
```

Output:

```
[#1] line 0 first writer
[#2] line 0 second writer
[#1] line 1 first writer
[#1] line 2 first writer
[#2] line 1 second writer
[#1] line 3 first writer
[#1] line 4 first writer
[#1] line 5 first writer
[#2] line 2 second writer
[#1] line 6 first writer
[#1] line 7 first writer
...
```

### For not safe writers:

``` GO
var writer bytes.Buffer
safeWriter := linebyline.NewSafeWriter(&writer)

wr1 := linebyline.NewByLineWriter(
    linebyline.WithOutWriter(os.Stdout),
)
defer wr1.Close()

wr2 := linebyline.NewByLineWriter(
    linebyline.WithOutWriter(os.Stdout),
)
defer wr2.Close()

fmt.Fprintln(wr1, "second writer")
fmt.Fprintln(wr2, "second writer")
// do something else ...

```

### Use flush function

``` GO
var writer bytes.Buffer
safeWriter := linebyline.NewSafeWriter(&writer)

wr1 := linebyline.NewByLineWriter(
    linebyline.WithFlushFunc(func(bytes []byte) error {
        print(string(bytes))
        return nil
    }),
)
defer wr1.Close()

wr2 := linebyline.NewByLineWriter(
    linebyline.WithFlushFunc(func(bytes []byte) error {
        print(string(bytes))
        return nil
    }),
)
defer wr2.Close()

fmt.Fprintln(wr1, "second writer")
fmt.Fprintln(wr2, "second writer")
// do something else ...
```

### Use custom end rune

``` GO
var writer bytes.Buffer
safeWriter := linebyline.NewSafeWriter(&writer)

wr1 := linebyline.NewByLineWriter(
    linebyline.WithOutWriter(os.Stdout),
    linebyline.WithEndRune('\t'),
)
defer wr1.Close()

wr2 := linebyline.NewByLineWriter(
    linebyline.WithOutWriter(os.Stdout),
    linebyline.WithEndRune('\t'),
)
defer wr2.Close()

fmt.Fprintln(wr1, "second writer\t")
fmt.Fprintln(wr2, "second writer\t")
// do something else ...
```