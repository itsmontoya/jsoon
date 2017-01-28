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

func appendString(val, s []byte) ([]byte, int) {
	var opened bool
	for i, b := range s {
		if b != charDoubleQuote {
			val = append(val, b)
			continue
		}

		if !opened {
			opened = true
			continue
		}

		return val, i
	}

	return val, -1
}

func appendNumber(val, s []byte) ([]byte, int) {
	for i, b := range s {
		if isNumber(b) {
			val = append(val, b)
			continue
		}

		switch b {
		case charComma, charSpace, charCloseCurly:
			return val, i - 1
		default:
			break
		}
	}

	return val, -1
}

func appendBool(val, s []byte) ([]byte, int) {
	switch s[0] {
	case charLowerT:
		if isTrueBytes(s[:4]) {
			return append(val, s[:4]...), 3
		}
	case charLowerF:
		if isFalseBytes(s[:5]) {
			return append(val, s[:5]...), 4
		}
	}

	return val, -1
}

func appendObject(val, s []byte) ([]byte, int) {
	var depth int
	for i, b := range s {
		val = append(val, b)

		switch b {
		case charCloseCurly:
			if depth--; depth == 0 {
				return val, i
			}
		case charOpenCurly:
			depth++
		}
	}

	return val, -1
}

func appendArray(val, s []byte) ([]byte, int) {
	var depth int
	for i, b := range s {
		val = append(val, b)

		switch b {
		case charCloseBracket:
			if depth--; depth == 0 {
				return val, i
			}
		case charOpenBracket:
			depth++
		}
	}

	return val, -1
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
