package linebyline

import (
	"io"
	"sync"
)

type WriterGroup struct {
	channel chan []byte
	writers []io.WriteCloser
	wg      *sync.WaitGroup
}

func NewWriterGroup(wr io.Writer) *WriterGroup {
	channel := make(chan []byte)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go writerDeamon(&wg, channel, wr)

	return &WriterGroup{
		channel: channel,
		writers: []io.WriteCloser{},
		wg:      &wg,
	}
}

func (wrg *WriterGroup) CreateWriter() io.WriteCloser {
	writer := byLineWriter{
		channel:         wrg.channel,
		trailingNewline: true,
	}

	wrg.writers = append(wrg.writers, &writer)

	return &writer
}

func (wrg *WriterGroup) Close() error {
	for _, v := range wrg.writers {
		v.Close()
	}

	close(wrg.channel)

	wrg.wg.Wait()

	return nil
}

func writerDeamon(wg *sync.WaitGroup, channel <-chan []byte, wr io.Writer) {
	defer wg.Done()
	for line := range channel {
		_, _ = wr.Write(line)
	}
}
