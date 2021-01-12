package buffer

import (
	"bytes"
)

type Buffer struct {
	*bytes.Buffer
}

func (p *Buffer) Close() error {
	p.Buffer.Reset()
	return nil
}

// New new buffer
func New(b *bytes.Buffer) *Buffer {
	if b == nil {
		b = bytes.NewBuffer(nil)
	}
	return &Buffer{b}
}
