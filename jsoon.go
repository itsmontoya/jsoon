package jsoon

import "github.com/missionMeteora/toolkit/errors"
import "strconv"

const (
	// ErrInvalidChar is returned when an invalid character is found within a provided json
	ErrInvalidChar = errors.Error("invalid character")
	// ErrUnexpectedEnd is returned when decoding data ends before the json is complete
	ErrUnexpectedEnd = errors.Error("input ended before expected")
	//ErrInvalidValue is returned when a Decodee is not present for an object or a ArrayDecodee is not present for an array
	ErrInvalidValue = errors.Error("invalid value provided")

	// ErrValueNotObject is returned when value is not an object
	ErrValueNotObject = errors.Error("value cannot be parsed as an object")
	// ErrValueNotArray is returned when value is not an array
	ErrValueNotArray = errors.Error("value cannot be parsed as an array")
	// ErrValueNotString is returned when value is not a string
	ErrValueNotString = errors.Error("value cannot be parsed as a string")
	// ErrValueNotBytes is returned when value is not bytes
	ErrValueNotBytes = errors.Error("value cannot be parsed as bytes")
	// ErrValueNotNumber is returned when value is not a number
	ErrValueNotNumber = errors.Error("value cannot be parsed as a number")
	// ErrValueNotBool is returned when value is not a boolean
	ErrValueNotBool = errors.Error("value cannot be parsed as a boolean")
)

const (
	valNil uint8 = iota
	valObject
	valArray
	valString
	valNumber
	valBool
)

var p = newPool(1024 * 4)

// Encodee is an item that has a Marshal helper func
type Encodee interface {
	MarshalJsoon(enc *Encoder) error
}

// ArrayEncodee is an item that has a Marshal helper func
type ArrayEncodee interface {
	MarshalJsoon(enc *ArrayEncoder) error
}

// Decodee is an item that has an Unmarshal helper func
type Decodee interface {
	UnmarshalJsoon(key string, val *Value) error
}

// ArrayDecodee is an array with decoding helper funcs
type ArrayDecodee interface {
	UnmarshalJsoon(val *Value) error
}

// Value represents a value
type Value struct {
	// Value type
	vt uint8
	// Reference decoder
	d *Decoder
}

// Object will associate a provided value with an object
func (v *Value) Object(val Decodee) (err error) {
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

	return strconv.ParseFloat(string(v.d.vb.Bytes()), 64)
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

// ReadByter is a byte reading interface
type ReadByter interface {
	Read([]byte) (int, error)
	ReadByte() (byte, error)
	UnreadByte() error
}
