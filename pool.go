package jsoon

import (
	"sync"
)

func newPool() *pool {
	var p pool
	p.p = sync.Pool{
		New: func() interface{} {
			return newBuffer()
		},
	}
	return &p
}

type pool struct {
	p sync.Pool
}

// Acquire will acquire a buffer from the pool
func (p *pool) Acquire() (buf *buffer) {
	var ok bool
	if buf, ok = p.p.Get().(*buffer); !ok {
		panic("invalid pool type")
	}

	return
}

// Release will release a buffer to the pool
func (p *pool) Release(buf *buffer) {
	buf.Reset()
	p.p.Put(buf)
}
