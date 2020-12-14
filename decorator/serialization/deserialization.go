package serialization

import (
	"encoding/gob"
	"io"

	"github.com/septemhill/translator/decorator"
)

type deserialization struct {
	d decorator.BufferDecorator
}

func (c *deserialization) Decorate() io.Reader {
	return c.d.Decorate()
}

func (c *deserialization) Decode(v interface{}) error {
	return gob.NewDecoder(c.d.Decorate()).Decode(v)
}

var _ decorator.BufferDecorator = (*deserialization)(nil)

func Deserialization(d decorator.BufferDecorator) *deserialization {
	return &deserialization{d: d}
}
