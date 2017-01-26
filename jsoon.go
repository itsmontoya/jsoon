package jsoon

import (
	"io"
	"sync"
)

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

// NewEncoder will return a new encoder
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// Encoder will handle the encoding
type Encoder struct {
	w io.Writer

	buf *buffer

	depth int
	child int
}

// Encode will marshal an Encodee
func (e *Encoder) Encode(value Encodee) {
	// Get parent's child value
	pc := e.child
	// Get parent's buffer
	pb := e.buf
	if pb != nil {
		e.w.Write(pb.Bytes())
		pb.Reset()
	}

	// Set child value to 0, since this is a new object
	e.child = 0
	// Increase depth
	e.depth++
	// Acquire buffer for this depth
	e.buf = acquireBuffer()

	e.buf.WriteByte('{')
	value.MarshalJsoon(e)
	e.buf.WriteByte('}')
	e.w.Write(e.buf.Bytes())

	// Release buffer for this depth
	releaseBuffer(e.buf)
	// Set buffer as the parent's buffer
	e.buf = pb
	// Reduce depth to the parent's level
	e.depth--
	// Set child value to parent's child value
	e.child = pc
}

// Object will marshal an Encodee
func (e *Encoder) Object(key string, value Encodee) {
	// Get parent's child value
	pc := e.child
	// Get parent's buffer
	pb := e.buf
	if pb != nil {
		e.w.Write(pb.Bytes())
		pb.Reset()
	}

	// Acquire buffer for this depth
	e.buf = acquireBuffer()

	if e.child > 0 {
		e.buf.WriteByte(',')
	}

	// Set child value to 0, since this is a new object
	e.child = 0
	// Increase depth
	e.depth++

	e.buf.WriteByte('"')
	e.buf.WriteString(key)
	e.buf.WriteString(`":{`)
	value.MarshalJsoon(e)
	e.buf.WriteByte('}')
	e.w.Write(e.buf.Bytes())

	// Release buffer for this depth
	releaseBuffer(e.buf)
	// Set buffer as the parent's buffer
	e.buf = pb
	// Reduce depth to the parent's level
	e.depth--
	// Set child value to parent's child value
	e.child = pc
}

// String will marshal a string
func (e *Encoder) String(key, value string) {
	if e.child > 0 {
		e.buf.WriteByte(',')
	}

	e.buf.WriteByte('"')
	e.buf.WriteString(key)

	e.buf.WriteString(`":"`)

	e.buf.WriteString(value)
	e.buf.WriteByte('"')

	e.child++
}

// Number will marshal a number
func (e *Encoder) Number(key string, value float64) {
	if e.child > 0 {
		e.buf.WriteByte(',')
	}

	e.buf.WriteByte('"')
	e.buf.WriteString(key)

	e.buf.WriteString(`":`)
	e.buf.WriteFloat64(value)
	e.child++
}

// Bool will marshal a boolean
func (e *Encoder) Bool(key string, value bool) {
	if e.child > 0 {
		e.buf.WriteByte(',')
	}

	e.buf.WriteByte('"')
	e.buf.WriteString(key)

	e.buf.WriteString(`":`)
	e.buf.WriteBool(value)
	e.child++
}

// Encodee is an item that has a Marshal helper func
type Encodee interface {
	MarshalJsoon(enc *Encoder) error
}
