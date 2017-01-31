package jsoon

import (
	"bytes"
)

var (
	trueBytes  = [4]byte{'t', 'r', 'u', 'e'}
	falseBytes = [5]byte{'f', 'a', 'l', 's', 'e'}
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
