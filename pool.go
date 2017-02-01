package jsoon

import "sync"

func newPool() *pool {
	var p pool
	p.p = sync.Pool{
		New: func() interface{} {
			return newBuffer()
		},
	}

	p.aep = sync.Pool{
		New: func() interface{} {
			return &ArrayEncoder{}
		},
	}

	return &p
}

type pool struct {
	p sync.Pool
	// array encoder pool
	aep sync.Pool
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

// AcquireAE will acquire an array encoder from the pool
func (p *pool) AcquireAE(e *Encoder) (ae *ArrayEncoder) {
	var ok bool
	if ae, ok = p.aep.Get().(*ArrayEncoder); !ok {
		panic("invalid pool type")
	}

	ae.e = e
	return
}

// ReleaseAE will release an array encoder to the pool
func (p *pool) ReleaseAE(ae *ArrayEncoder) {
	ae.e = nil
	p.aep.Put(ae)
}
