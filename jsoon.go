package jsoon

import "errors"

var (
	// ErrInvalidChar is returned when an invalid character is found within a provided json
	ErrInvalidChar = errors.New("invalid character")
	// ErrUnexpectedEnd is returned when decoding data ends before the json is complete
	ErrUnexpectedEnd = errors.New("input ended before expected")
	//ErrInvalidValue is returned when a Decodee is not present for an object or a ArrayDecodee is not present for an array
	ErrInvalidValue = errors.New("invalid value provided")

	// ErrValueNotObject is returned when value is not an object
	ErrValueNotObject = errors.New("value cannot be parsed as an object")
	// ErrValueNotArray is returned when value is not an array
	ErrValueNotArray = errors.New("value cannot be parsed as an array")
	// ErrValueNotString is returned when value is not a string
	ErrValueNotString = errors.New("value cannot be parsed as a string")
	// ErrValueNotBytes is returned when value is not bytes
	ErrValueNotBytes = errors.New("value cannot be parsed as bytes")
	// ErrValueNotNumber is returned when value is not a number
	ErrValueNotNumber = errors.New("value cannot be parsed as a number")
	// ErrValueNotBool is returned when value is not a boolean
	ErrValueNotBool = errors.New("value cannot be parsed as a boolean")
)

const (
	valNil uint8 = iota
	valObject
	valArray
	valString
	valNumber
	valBool
)

var p = newPool()

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

// ReadByter is a byte reading interface
type ReadByter interface {
	Read([]byte) (int, error)
	ReadByte() (byte, error)
	UnreadByte() error
}
