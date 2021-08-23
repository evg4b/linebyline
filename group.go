package linebyline

import (
	"io"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type WriterGroup struct {
	writers []io.WriteCloser
	mu      sync.Mutex
	wr      io.Writer
}

func NewWriterGroup(wr io.Writer) *WriterGroup {
	return &WriterGroup{
		writers: []io.WriteCloser{},
		wr:      wr,
		mu:      sync.Mutex{},
	}
}

func (wrg *WriterGroup) CreateWriter() io.WriteCloser {
	writer := byLineWriter{
		originalWriter:  wrg.wr,
		mu:              &wrg.mu,
		trailingNewline: true,
	}

	wrg.writers = append(wrg.writers, &writer)

	return &writer
}

func (wrg *WriterGroup) Close() error {
	var errors *multierror.Error
	for _, v := range wrg.writers {
		err := v.Close()
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	return errors.ErrorOrNil()
}
