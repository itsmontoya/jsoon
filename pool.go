package jsoon

import (
	"github.com/itsmontoya/mailbox"
)

func newPool(sz int) *pool {
	var p pool
	p.buf = mailbox.New(sz)
	p.val = mailbox.New(sz)
	return &p
}

type pool struct {
	buf *mailbox.Mailbox
	val *mailbox.Mailbox
}

// acquireBuffer will acquire a buffer from the pool
func (p *pool) acquireBuffer() (buf *buffer) {
	var (
		v  interface{}
		sc mailbox.StateCode
		ok bool
	)

	if v, sc = p.buf.Receive(false); sc == mailbox.StateEmpty {
		return newBuffer()
	}

	if buf, ok = v.(*buffer); !ok {
		panic("invalid pool type")
	}

	return
}

// releaseBuffer will release a buffer to the pool
func (p *pool) releaseBuffer(buf *buffer) {
	buf.Reset()
	// If the mailbox is full, move along
	p.buf.Send(buf, false)
}
