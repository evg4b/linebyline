package linebyline

import (
	"bytes"
	"io"
	"sync"
)

type byLineWriter struct {
	buffer         bytes.Buffer
	mu             *sync.Mutex
	originalWriter io.Writer
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

		wr.mu.Lock()
		defer wr.mu.Unlock()

		_, err = wr.originalWriter.Write(wr.buffer.Bytes())

		return err
	}

	return nil
}
