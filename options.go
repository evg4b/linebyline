package linebyline

import "io"

type Option func(*byLineWriter)

func WithOutWriter(w io.Writer) Option {
	return func(wr *byLineWriter) {
		wr.flushFn = func(bytes []byte) error {
			_, err := w.Write(bytes)
			return err
		}
	}
}

func WithFlushFunc(fn func([]byte) error) Option {
	return func(wr *byLineWriter) {
		wr.flushFn = fn
	}
}

func WithEndRune(endRune rune) Option {
	return func(wr *byLineWriter) {
		wr.endRune = byte(endRune)
	}
}
