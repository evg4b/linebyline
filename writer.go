package linebyline

import (
	"bytes"
	"io"
)

type byLineWriter struct {
	buffer bytes.Buffer
	out    io.Writer
}

func NewByLineWriter(out io.Writer) *byLineWriter {
	return &byLineWriter{
		buffer: bytes.Buffer{},
		out:    out,
	}
}

func (wr *byLineWriter) Close() error {
	return wr.flush()
}

func (wr *byLineWriter) Write(payload []byte) (int, error) {
	for _, b := range payload {
		if b == '\n' {
			err := wr.flush()
			if err != nil {
				return 0, err
			}
		} else {
			err := wr.buffer.WriteByte(b)
			if err != nil {
				return 0, err
			}
		}
	}

	return len(payload), nil
}

func (wr *byLineWriter) flush() error {
	defer wr.buffer.Reset()

	if wr.buffer.Len() > 0 {
		_, err := wr.buffer.WriteRune('\n')
		if err != nil {
			return err
		}

		_, err = io.Copy(wr.out, &wr.buffer)

		return err
	}

	return nil
}
