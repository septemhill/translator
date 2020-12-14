package compress

import (
	"compress/gzip"
	"io"

	"github.com/septemhill/translator/decorator"
)

type compress struct {
	d decorator.BufferDecorator
}

func (c *compress) Decorate() io.Reader {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		gw := gzip.NewWriter(pw)
		defer gw.Close()
		io.Copy(gw, c.d.Decorate())
	}()
	return pr
}

var _ decorator.BufferDecorator = (*compress)(nil)

func Compress(d decorator.BufferDecorator) *compress {
	return &compress{d: d}
}
