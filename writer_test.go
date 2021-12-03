package linebyline_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	. "github.com/evg4b/linebyline"
	"github.com/stretchr/testify/assert"
)

func Test_ByLineWriter_BasicUsage(t *testing.T) {
	var writer bytes.Buffer

	wr1 := NewByLineWriter(&writer)
	wr2 := NewByLineWriter(&writer)

	fmt.Fprint(wr1, "This is first")
	fmt.Fprint(wr2, "This is second")
	fmt.Fprintln(wr1, " writer")
	fmt.Fprintln(wr2, " writer")

	assert.NoError(t, wr1.Close())
	assert.NoError(t, wr2.Close())
	assert.Equal(t, "This is first writer\nThis is second writer\n", writer.String())
}

func Test_ByLineWriter_InterruptedWriting(t *testing.T) {
	var writer bytes.Buffer

	wr1 := NewByLineWriter(&writer)
	wr2 := NewByLineWriter(&writer)

	fmt.Fprint(wr1, "This is first")
	fmt.Fprint(wr2, "This is second")
	fmt.Fprint(wr1, " writer")
	fmt.Fprint(wr2, " writer")

	assert.NoError(t, wr1.Close())
	assert.NoError(t, wr2.Close())
	assert.Equal(t, "This is first writer\nThis is second writer\n", writer.String())
}

func Test_ByLineWriter_NolLockWriter(t *testing.T) {
	var writer bytes.Buffer
	safeWriter := NewSafeWriter(&writer)

	wr1 := NewByLineWriter(safeWriter)
	wr2 := NewByLineWriter(safeWriter)

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
}
