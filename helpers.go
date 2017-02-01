package jsoon

import (
	"bytes"
	"unsafe"
)

var (
	trueBytes  = [4]byte{'t', 'r', 'u', 'e'}
	falseBytes = [5]byte{'f', 'a', 'l', 's', 'e'}
	nullBytes  = [4]byte{'n', 'u', 'l', 'l'}
)

func isTrueBytes(s []byte) bool {
	return bytes.Compare(trueBytes[:], s) == 0
}

func isFalseBytes(s []byte) bool {
	return bytes.Compare(falseBytes[:], s) == 0
}

func isLetter(b byte) bool {
	if b >= charLowerA && b <= charLowerZ {
		return true
	}

	if b >= charUpperA && b <= charUpperZ {
		return true
	}

	return false
}

func isNumber(b byte) bool {
	return b >= charZero && b <= charNine
}

func isWhitespace(b byte) bool {
	return b == charSpace || b == charTab || b == charNewline
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
