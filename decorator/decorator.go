package decorator

import (
	"bytes"
	"io"
)

type BufferDecorator interface {
	Decorate() io.Reader
}

type bufferDecorator struct {
	r io.Reader
}

func (d *bufferDecorator) Decorate() io.Reader {
	return d.r
}

func NewBufferDecorator(b []byte) *bufferDecorator {
	return &bufferDecorator{
		r: bytes.NewBuffer(b),
	}
}
