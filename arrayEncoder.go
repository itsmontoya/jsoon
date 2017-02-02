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
func (a *ArrayEncoder) Object(value Encodee) (err error) {
	if a.e.depth == 0 {
		// Acquire buffer for this depth
		a.e.buf = p.Acquire()
	}

	// Increase depth
	a.e.depth++

	// Get parent's child value
	pc := a.e.child

	if a.e.child > 0 {
		a.e.buf.WriteByte(charComma)
	}

	// Set child value to 0, since this is a new object
	a.e.child = 0

	a.e.buf.WriteByte(charOpenCurly)

	if err = value.MarshalJsoon(a.e); err != nil {
		return
	}

	a.e.buf.WriteByte(charCloseCurly)

	if _, err = a.e.w.Write(a.e.buf.Bytes()); err != nil {
		return
	}

	a.e.buf.Reset()

	// Set child value to parent's child value
	a.e.child = pc
	a.e.child++

	// Reduce depth to the parent's level
	a.e.depth--
	return
}

// Array will marshal an array
func (a *ArrayEncoder) Array(value ArrayEncodee) (err error) {
	if a.e.depth == 0 {
		// Acquire buffer for this depth
		a.e.buf = p.Acquire()
	}

	// Increase depth
	a.e.depth++

	// Get parent's child value
	pc := a.e.child

	if a.e.child > 0 {
		a.e.buf.WriteByte(charComma)
	}

	// Set child value to 0, since this is a new object
	a.e.child = 0

	a.e.buf.WriteByte(charOpenBracket)

	if err = value.MarshalJsoon(a); err != nil {
		return
	}

	a.e.buf.WriteByte(charCloseBracket)

	if _, err = a.e.w.Write(a.e.buf.Bytes()); err != nil {
		return
	}

	a.e.buf.Reset()

	// Set child value to parent's child value
	a.e.child = pc
	a.e.child++

	// Reduce depth to the parent's level
	a.e.depth--
	return
}

// String will marshal a string
func (a *ArrayEncoder) String(value string) {
	if a.e.child > 0 {
		a.e.buf.WriteByte(charComma)
	}

	a.e.buf.WriteByte(charDoubleQuote)
	a.e.buf.WriteString(value)
	a.e.buf.WriteByte(charDoubleQuote)
	a.e.child++
}

// Number will marshal a number
func (a *ArrayEncoder) Number(value float64) {
	if a.e.child > 0 {
		a.e.buf.WriteByte(charComma)
	}

	a.e.buf.WriteFloat64(value)
	a.e.child++
}

// Bool will marshal a boolean
func (a *ArrayEncoder) Bool(value bool) {
	if a.e.child > 0 {
		a.e.buf.WriteByte(charComma)
	}

	a.e.buf.WriteBool(value)
	a.e.child++
}
