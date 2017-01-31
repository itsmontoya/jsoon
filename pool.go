package jsoon

import (
	"github.com/itsmontoya/mailbox"
)

func newPool(sz int) *pool {
	var p pool
	p.mb = mailbox.New(sz)
	return &p
}

type pool struct {
	mb *mailbox.Mailbox
}

// Acquire will acquire a buffer from the pool
func (p *pool) Acquire() (buf *buffer) {
	var (
		v  interface{}
		sc mailbox.StateCode
		ok bool
	)

	if v, sc = p.mb.Receive(false); sc == mailbox.StateEmpty {
		return newBuffer()
	}

	if buf, ok = v.(*buffer); !ok {
		panic("invalid pool type")
	}

	return
}

// Release will release a buffer to the pool
func (p *pool) Release(buf *buffer) {
	buf.Reset()
	// If the mailbox is full, move along
	p.mb.Send(buf, false)
}