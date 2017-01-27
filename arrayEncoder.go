package jsoon

func newArrayEncoder(e *Encoder) *ArrayEncoder {
	return &ArrayEncoder{e: e}
}

// ArrayEncoder is an encoder for arrays
// TODO: Use more helper funcs to reduce code duplication between Encoder and  ArrayEncoder fn's
type ArrayEncoder struct {
	e *Encoder
}

// Object will marshal an Encodee
func (a *ArrayEncoder) Object(value Encodee) {
	// Get parent's child value
	pc := a.e.child
	// Get parent's buffer
	pb := a.e.buf
	if pb != nil {
		a.e.w.Write(pb.Bytes())
		pb.Reset()
	}

	// Acquire buffer for this depth
	a.e.buf = acquireBuffer()

	if a.e.child > 0 {
		a.e.buf.WriteByte(',')
	}

	// Set child value to 0, since this is a new object
	a.e.child = 0
	// Increase depth
	a.e.depth++

	a.e.buf.WriteByte('{')
	value.MarshalJsoon(a.e)
	a.e.buf.WriteByte('}')
	a.e.w.Write(a.e.buf.Bytes())

	// Release buffer for this depth
	releaseBuffer(a.e.buf)
	// Set buffer as the parent's buffer
	a.e.buf = pb
	// Reduce depth to the parent's level
	a.e.depth--
	// Set child value to parent's child value
	a.e.child = pc
	a.e.child++
}

// Array will marshal an array
func (a *ArrayEncoder) Array(value ArrayEncodee) {
	// Get parent's child value
	pc := a.e.child
	// Get parent's buffer
	pb := a.e.buf
	if pb != nil {
		a.e.w.Write(pb.Bytes())
		pb.Reset()
	}

	// Acquire buffer for this depth
	a.e.buf = acquireBuffer()

	if a.e.child > 0 {
		a.e.buf.WriteByte(',')
	}

	// Set child value to 0, since this is a new object
	a.e.child = 0
	// Increase depth
	a.e.depth++

	a.e.buf.WriteByte('[')
	value.MarshalJsoon(a)
	a.e.buf.WriteByte(']')
	a.e.w.Write(a.e.buf.Bytes())

	// Release buffer for this depth
	releaseBuffer(a.e.buf)
	// Set buffer as the parent's buffer
	a.e.buf = pb
	// Reduce depth to the parent's level
	a.e.depth--
	// Set child value to parent's child value
	a.e.child = pc
	a.e.child++
}

// String will marshal a string
func (a *ArrayEncoder) String(value string) {
	if a.e.child > 0 {
		a.e.buf.WriteByte(',')
	}

	a.e.buf.WriteByte('"')
	a.e.buf.WriteString(value)
	a.e.buf.WriteByte('"')
	a.e.child++
}

// Number will marshal a number
func (a *ArrayEncoder) Number(value float64) {
	if a.e.child > 0 {
		a.e.buf.WriteByte(',')
	}

	a.e.buf.WriteFloat64(value)
	a.e.child++
}

// Bool will marshal a boolean
func (a *ArrayEncoder) Bool(value bool) {
	if a.e.child > 0 {
		a.e.buf.WriteByte(',')
	}

	a.e.buf.WriteBool(value)
	a.e.child++
}
