package linebyline

import (
	"io"
	"sync"
)

type safeWriter struct {
	mu  sync.Mutex
	out io.Writer
}

func NewSafeWriter(out io.Writer) io.Writer {
	return &safeWriter{
		mu:  sync.Mutex{},
		out: out,
	}
}

func (wr *safeWriter) Write(payload []byte) (int, error) {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	return wr.out.Write(payload)
}
