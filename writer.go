package linebyline

import (
	"bytes"
	"io"
)

type byLineWriter struct {
	buffer  bytes.Buffer
	endRune byte
	flushFn func([]byte) error
}

func NewByLineWriter(options ...Option) io.WriteCloser {
	rw := &byLineWriter{
		buffer:  bytes.Buffer{},
		endRune: '\n',
		flushFn: func(bytes []byte) error {
			return nil
		},
	}

	for _, option := range options {
		option(rw)
	}

	return rw
}

func (wr *byLineWriter) Close() error {
	return wr.flush()
}

func (wr *byLineWriter) Write(payload []byte) (int, error) {
	for _, b := range payload {
		if err := wr.writeByte(b); err != nil {
			return 0, err
		}
	}

	return len(payload), nil
}

func (wr *byLineWriter) writeByte(b byte) error {
	if b == wr.endRune {
		return wr.flush()
	}

	return wr.buffer.WriteByte(b)
}

func (wr *byLineWriter) flush() error {
	defer wr.buffer.Reset()

	if wr.buffer.Len() > 0 {
		if err := wr.buffer.WriteByte(wr.endRune); err != nil {
			return err
		}

		return wr.flushFn(wr.buffer.Bytes())
	}

	return nil
}
