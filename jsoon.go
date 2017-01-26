package jsoon

import "sync"

var p = sync.Pool{
	New: func() interface{} {
		return newBuffer()
	},
}

func acquireBuffer() (buf *buffer) {
	var ok bool
	if buf, ok = p.Get().(*buffer); !ok {
		panic("invalid pool type")
	}

	return
}

func releaseBuffer(buf *buffer) {
	buf.Reset()
	p.Put(buf)
}

// Encodee is an item that has a Marshal helper func
type Encodee interface {
	MarshalJsoon(enc *Encoder) error
}
