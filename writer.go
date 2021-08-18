package linebyline

import (
	"bytes"
)

type byLineWriter struct {
	buffer          bytes.Buffer
	channel         chan<- []byte
	trailingNewline bool
}

func (wr *byLineWriter) Close() error {
	payload := wr.buffer.Bytes()
	if len(payload) > 0 {
		wr.flush()
	}

	return nil
}

func (wr *byLineWriter) Write(payload []byte) (int, error) {
	for _, b := range payload {
		if wr.trailingNewline {
			wr.trailingNewline = false
		}

		if b == '\n' {
			wr.trailingNewline = true
			wr.flush()
		} else {
			wr.buffer.WriteByte(b)
		}
	}

	return len(payload), nil
}

func (wr *byLineWriter) flush() {
	_, _ = wr.buffer.WriteRune('\n')
	wr.channel <- wr.buffer.Bytes()
	wr.buffer.Reset()
}
