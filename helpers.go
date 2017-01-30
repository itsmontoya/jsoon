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

func (d *Decoder) appendString(val []byte) ([]byte, error) {
	for b, err := d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if b != charDoubleQuote {
			val = append(val, b)
			continue
		}

		return val, nil
	}

	return nil, ErrInvalidChar
}

func (d *Decoder) appendNumber(val []byte) ([]byte, error) {
	for b, err := d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if isNumber(b) {
			val = append(val, b)
			continue
		}

		switch b {
		case charComma, charSpace, charCloseCurly:
			d.r.UnreadByte()
			return val, nil
		default:
			break
		}
	}

	return nil, ErrInvalidChar
}

func (d *Decoder) appendTrue(val []byte) ([]byte, error) {
	var i int
	for b, err := d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if b != trueBytes[i] {
			return val, ErrInvalidChar
		}

		val = append(val, b)

		if i++; i == 4 {
			break
		}
	}

	return val, nil
}

func (d *Decoder) appendFalse(val []byte) ([]byte, error) {
	var i int
	for b, err := d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if b != falseBytes[i] {
			return nil, ErrInvalidChar
		}

		val = append(val, b)

		if i++; i == 5 {
			break
		}
	}

	return nil, nil
}

func (d *Decoder) appendObject(val []byte) ([]byte, error) {
	var depth int
	for b, err := d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		val = append(val, b)
		switch b {
		case charCloseCurly:
			if depth--; depth == 0 {
				return val, nil
			}

		case charOpenCurly:
			depth++
		}
	}

	return nil, ErrInvalidChar
}

func (d *Decoder) appendArray(val []byte) ([]byte, error) {
	var depth int
	for b, err := d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		val = append(val, b)
		switch b {
		case charCloseBracket:
			if depth--; depth == 0 {
				return val, nil
			}
		case charOpenBracket:
			depth++
		}
	}

	return nil, ErrInvalidChar
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
