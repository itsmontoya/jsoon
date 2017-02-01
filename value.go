package jsoon

import "strconv"

// Value represents a value
type Value struct {
	// Value type
	vt uint8
	// Reference decoder
	d *Decoder
}

// Object will associate a provided value with an object
func (v *Value) Object(val Decodee) (err error) {
	if v.vt == valNil {
		return
	}

	if v.vt != valObject {
		return ErrValueNotObject
	}

	v.d.kb.Reset()
	if err = v.d.decodeObject(val); err != nil {
		return
	}

	return
}

// Array will associate a provided value with an array
func (v *Value) Array(val ArrayDecodee) (err error) {
	if v.vt == valNil {
		return
	}

	if v.vt != valArray {
		return ErrValueNotArray
	}

	v.d.kb.Reset()
	if err = v.d.decodeArray(val); err != nil {
		return
	}

	return
}

func (v *Value) String() (val string, err error) {
	if v.vt != valString {
		err = ErrValueNotString
		return
	}

	val = string(v.d.vb.Bytes())
	return
}

// Bytes will return bytes
// Note: Please do not hold onto the val after it's initially returned. Use it, let it spread it's wings and fly
func (v *Value) Bytes() (val []byte, err error) {
	if v.vt != valString {
		err = ErrValueNotBytes
		return
	}

	val = v.d.vb.Bytes()
	return
}

// Number will return a number value
func (v *Value) Number() (val float64, err error) {
	if v.vt != valNumber {
		err = ErrValueNotNumber
		return
	}

	return strconv.ParseFloat(unsafeString(v.d.vb.Bytes()), 64)
}

// Bool will return a boolean value
func (v *Value) Bool() (val bool, err error) {
	if v.vt != valBool {
		err = ErrValueNotBool
		return
	}

	val = len(v.d.vb.Bytes()) == 4
	return
}
