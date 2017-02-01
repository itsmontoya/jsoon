package jsoon

import "io"

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
	if e.depth == 0 {
		// Acquire buffer for this depth
		e.buf = p.Acquire()
	}

	// Increase depth
	e.depth++

	// Get parent's child value
	pc := e.child

	// Set child value to 0, since this is a new object
	e.child = 0

	e.buf.WriteByte(charOpenCurly)
	value.MarshalJsoon(e)
	e.buf.WriteByte(charCloseCurly)
	e.w.Write(e.buf.Bytes())
	e.buf.Reset()

	// Set child value to parent's child value
	e.child = pc

	// Reduce depth to the parent's level
	e.depth--

	if e.depth == 0 {
		p.Release(e.buf)
		e.buf = nil
	}
}

// Object will marshal an Encodee
func (e *Encoder) Object(key string, value Encodee) {
	if e.depth == 0 {
		// Acquire buffer for this depth
		e.buf = p.Acquire()
	}

	// Increase depth
	e.depth++

	// Get parent's child value
	pc := e.child

	if e.child > 0 {
		e.buf.WriteByte(charComma)
	}

	// Set child value to 0, since this is a new object
	e.child = 0

	e.buf.WriteByte(charDoubleQuote)
	e.buf.WriteString(key)
	e.buf.WriteString(`":{`)
	value.MarshalJsoon(e)
	e.buf.WriteByte(charCloseCurly)
	e.w.Write(e.buf.Bytes())
	e.buf.Reset()

	// Reduce depth to the parent's level
	e.depth--

	// Set child value to parent's child value
	e.child = pc
	e.child++
}

// Array will marshal an array
func (e *Encoder) Array(key string, value ArrayEncodee) {
	if e.depth == 0 {
		// Acquire buffer for this depth
		e.buf = p.Acquire()
	}

	// Increase depth
	e.depth++

	// Get parent's child value
	pc := e.child

	if pc > 0 {
		e.buf.WriteByte(charComma)
	}

	// Set child value to 0, since this is a new object
	e.child = 0

	e.buf.WriteByte(charDoubleQuote)
	e.buf.WriteString(key)
	e.buf.WriteString(`":[`)
	ae := p.AcquireAE(e)
	value.MarshalJsoon(ae)
	p.ReleaseAE(ae)
	e.buf.WriteByte(charCloseBracket)
	e.w.Write(e.buf.Bytes())
	e.buf.Reset()

	// Reduce depth to the parent's level
	e.depth--

	// Set child value to parent's child value
	e.child = pc
	e.child++
}

// String will escape and marshal a string
func (e *Encoder) String(key, value string) {
	if e.child > 0 {
		e.buf.WriteByte(charComma)
	}

	e.buf.WriteByte(charDoubleQuote)
	e.buf.WriteString(key)

	e.buf.WriteString(`":"`)

	e.buf.WriteEscapedString(value)
	e.buf.WriteByte(charDoubleQuote)

	e.child++
}

// UnsafeString will will marshal a string without escaping
// Note: Only use this if you are CERTAIN that your value does not contain any quotes
func (e *Encoder) UnsafeString(key, value string) {
	if e.child > 0 {
		e.buf.WriteByte(charComma)
	}

	e.buf.WriteByte(charDoubleQuote)
	e.buf.WriteString(key)

	e.buf.WriteString(`":"`)

	e.buf.WriteString(value)
	e.buf.WriteByte(charDoubleQuote)

	e.child++
}

// Number will marshal a number
func (e *Encoder) Number(key string, value float64) {
	if e.child > 0 {
		e.buf.WriteByte(charComma)
	}

	e.buf.WriteByte(charDoubleQuote)
	e.buf.WriteString(key)

	e.buf.WriteString(`":`)
	e.buf.WriteFloat64(value)
	e.child++
}

// Bool will marshal a boolean
func (e *Encoder) Bool(key string, value bool) {
	if e.child > 0 {
		e.buf.WriteByte(charComma)
	}

	e.buf.WriteByte(charDoubleQuote)
	e.buf.WriteString(key)

	e.buf.WriteString(`":`)
	e.buf.WriteBool(value)
	e.child++
}
