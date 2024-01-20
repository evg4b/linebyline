package linebyline_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/evg4b/linebyline"
	"github.com/stretchr/testify/assert"
)

func TestByLineWriter(t *testing.T) {
	t.Run("Basic usage", func(t *testing.T) {
		var writer bytes.Buffer

		wr1 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(&writer),
		)
		wr2 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(&writer),
		)

		fmt.Fprint(wr1, "This is first")
		fmt.Fprint(wr2, "This is second")
		fmt.Fprintln(wr1, " writer")
		fmt.Fprintln(wr2, " writer")

		assert.NoError(t, wr1.Close())
		assert.NoError(t, wr2.Close())
		assert.Equal(t, "This is first writer\nThis is second writer\n", writer.String())
	})

	t.Run("Interrupted writing", func(t *testing.T) {
		var writer bytes.Buffer

		wr1 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(&writer),
		)
		wr2 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(&writer),
		)

		fmt.Fprint(wr1, "This is first")
		fmt.Fprint(wr2, "This is second")
		fmt.Fprint(wr1, " writer")
		fmt.Fprint(wr2, " writer")

		assert.NoError(t, wr1.Close())
		assert.NoError(t, wr2.Close())
		assert.Equal(t, "This is first writer\nThis is second writer\n", writer.String())
	})

	t.Run("No lock writer", func(t *testing.T) {
		var writer bytes.Buffer
		safeWriter := linebyline.NewSafeWriter(&writer)

		wr1 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(safeWriter),
		)
		wr2 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(safeWriter),
		)

		wg := sync.WaitGroup{}

		fmt.Fprint(wr1, "wr1")

		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, data := range []string{"wr2-data1", "wr2-data2", "wr2-data3"} {
				fmt.Fprintln(wr2, data)
			}
		}()

		wg.Wait()

		assert.NoError(t, wr1.Close())
		assert.NoError(t, wr2.Close())
		assert.Equal(t, "wr2-data1\nwr2-data2\nwr2-data3\nwr1\n", writer.String())
	})

	t.Run("Custom end rune", func(t *testing.T) {
		var writer bytes.Buffer

		wr1 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(&writer),
			linebyline.WithEndRune('\t'),
		)
		defer wr1.Close()

		wr2 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(&writer),
			linebyline.WithEndRune('\t'),
		)
		defer wr1.Close()

		fmt.Fprint(wr1, "This is first")
		fmt.Fprint(wr2, "This is second")
		fmt.Fprint(wr1, " writer\t")
		fmt.Fprint(wr2, " writer\t")

		assert.Equal(t, "This is first writer\tThis is second writer\t", writer.String())
	})

	t.Run("Custom flush function", func(t *testing.T) {
		var items []string

		wr1 := linebyline.NewByLineWriter(
			linebyline.WithFlushFunc(func(bytes []byte) error {
				items = append(items, string(bytes))
				return nil
			}),
		)
		wr2 := linebyline.NewByLineWriter(
			linebyline.WithFlushFunc(func(bytes []byte) error {
				items = append(items, string(bytes))
				return nil
			}),
		)

		wg := sync.WaitGroup{}

		fmt.Fprint(wr1, "wr1")

		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, data := range []string{"wr2-data1", "wr2-data2", "wr2-data3"} {
				fmt.Fprintln(wr2, data)
			}
		}()

		wg.Wait()

		assert.NoError(t, wr1.Close())
		assert.NoError(t, wr2.Close())
		assert.Equal(t, []string{"wr2-data1\n", "wr2-data2\n", "wr2-data3\n", "wr1\n"}, items)
	})

	t.Run("Omit new line rune", func(t *testing.T) {
		var writer bytes.Buffer

		wr1 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(&writer),
			linebyline.OmitNewLineRune(),
		)
		wr2 := linebyline.NewByLineWriter(
			linebyline.WithOutWriter(&writer),
			linebyline.OmitNewLineRune(),
		)

		fmt.Fprint(wr1, "This is first")
		fmt.Fprint(wr2, "This is second")
		fmt.Fprint(wr1, " writer")
		fmt.Fprint(wr2, " writer")

		assert.NoError(t, wr1.Close())
		assert.NoError(t, wr2.Close())
		assert.Equal(t, "This is first writerThis is second writer", writer.String())
	})
}
