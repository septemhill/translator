package serialization

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/septemhill/translator/decorator"
	"github.com/septemhill/translator/service"
)

type serialization struct {
	tr *service.Translation
}

func (c *serialization) Decorate() io.Reader {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(c.tr)
	return &b
}

var _ decorator.BufferDecorator = (*serialization)(nil)

func Serialization(tr *service.Translation) *serialization {
	return &serialization{tr: tr}
}
