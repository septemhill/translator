package compress

import (
	"compress/gzip"
	"io"

	"github.com/septemhill/translator/decorator"
)

type decompress struct {
	d decorator.BufferDecorator
}

func (c *decompress) Decorate() io.Reader {
	r, _ := gzip.NewReader(c.d.Decorate())
	return r
}

var _ decorator.BufferDecorator = (*decompress)(nil)

func Decompress(d decorator.BufferDecorator) *decompress {
	return &decompress{d: d}
}
