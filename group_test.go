package linebyline_test

import (
	"bytes"
	"fmt"
	"linebyline"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriterGroup_BasicUsage(t *testing.T) {
	var writer bytes.Buffer

	group := linebyline.NewWriterGroup(&writer)
	wr1 := group.CreateWriter()
	wr2 := group.CreateWriter()

	fmt.Fprint(wr1, "This is first")
	fmt.Fprint(wr2, "This is second")
	fmt.Fprintln(wr1, " writer")
	fmt.Fprintln(wr2, " writer")

	wr1.Close()
	wr2.Close()

	assert.NoError(t, group.Close())
	assert.Equal(t, "This is first writer\nThis is second writer\n", writer.String())
}

func Test_WriterGroup_InterruptedWriting(t *testing.T) {
	var writer bytes.Buffer

	group := linebyline.NewWriterGroup(&writer)
	wr1 := group.CreateWriter()
	wr2 := group.CreateWriter()

	fmt.Fprint(wr1, "This is first")
	fmt.Fprint(wr2, "This is second")
	fmt.Fprint(wr1, " writer")
	fmt.Fprint(wr2, " writer")

	wr1.Close()
	wr2.Close()

	assert.NoError(t, group.Close())
	assert.Equal(t, "This is first writer\nThis is second writer\n", writer.String())
}
func Test_WriterGroup_ClosedGroup(t *testing.T) {
	var writer bytes.Buffer

	group := linebyline.NewWriterGroup(&writer)
	wr1 := group.CreateWriter()
	wr2 := group.CreateWriter()

	fmt.Fprint(wr1, "This is first")
	fmt.Fprint(wr2, "This is second")
	fmt.Fprint(wr1, " writer")
	fmt.Fprint(wr2, " writer")

	assert.NoError(t, group.Close())
	assert.Equal(t, "This is first writer\nThis is second writer\n", writer.String())
}
