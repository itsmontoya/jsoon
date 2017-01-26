package jsoon

import (
	"strconv"
)

func newBuffer() *buffer {
	return &buffer{
		s: make([]byte, 0, 32),
	}
}

type buffer struct {
	s []byte
}

func (b *buffer) Write(v []byte) {
	b.s = append(b.s, v...)
}

func (b *buffer) WriteByte(v byte) {
	b.s = append(b.s, v)
}

func (b *buffer) WriteString(v string) {
	b.s = append(b.s, v...)
}

func (b *buffer) WriteFloat64(v float64) {
	b.s = strconv.AppendFloat(b.s, v, 'f', -1, 64)
}

func (b *buffer) WriteBool(v bool) {
	b.s = strconv.AppendBool(b.s, v)
}

func (b *buffer) Bytes() []byte {
	return b.s
}

func (b *buffer) Reset() {
	b.s = b.s[:0]
}
