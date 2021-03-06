package jsoon

import (
	"bufio"
	"io"
)

const (
	osStart uint8 = iota
	osKey
	osPreSeparator
	osValue
	osPostValue
	osEnd
)

const (
	asStart uint8 = iota
	asValue
	asPostValue
	asEnd
)

const (
	vsStart uint8 = iota
	vsOpen
	vsEnd
)

const (
	charSpace        = ' '
	charTab          = '\t'
	charNewline      = '\n'
	charDoubleQuote  = '"'
	charSingleQuote  = '\''
	charBackslash    = '\\'
	charOpenCurly    = '{'
	charCloseCurly   = '}'
	charOpenBracket  = '['
	charCloseBracket = ']'
	charColon        = ':'
	charComma        = ','
	charZero         = '0'
	charNine         = '9'
	charLowerA       = 'a'
	charLowerN       = 'n'
	charLowerZ       = 'z'
	charUpperA       = 'A'
	charUpperZ       = 'Z'
	charLowerT       = 't'
	charLowerF       = 'f'
	charHyphen       = '-'
)

// NewDecoder will return a new Decoder
func NewDecoder(r io.Reader) *Decoder {
	var (
		d  Decoder
		ok bool
	)

	if d.r, ok = r.(ReadByter); !ok {
		d.r = bufio.NewReader(r)
	}

	d.v.d = &d
	return &d
}

// Decoder handles decoding
type Decoder struct {
	r ReadByter
	// key buffer
	kb *buffer
	// value buffer
	vb *buffer
	// decode count
	dc int

	v Value
}

// Decode will decode
func (d *Decoder) Decode(value interface{}) (err error) {
	var b byte

	if d.dc == 0 {
		d.kb = p.Acquire()
		d.vb = p.Acquire()
	}
	d.dc++

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if isWhitespace(b) {
			continue
		}

		switch b {
		case charOpenCurly:
			dec, ok := value.(Decodee)
			if !ok {
				err = ErrInvalidValue
				goto END
			}

			err = d.decodeObject(dec)
			goto END

		case charOpenBracket:
			dec, ok := value.(ArrayDecodee)
			if !ok {
				err = ErrInvalidChar
				goto END
			}

			err = d.decodeArray(dec)
			goto END

		default:
			err = ErrInvalidChar
			goto END
		}
	}

END:
	d.dc--
	if d.dc == 0 {
		p.Release(d.kb)
		p.Release(d.vb)
		d.kb = nil
		d.vb = nil
	}
	return
}

func (d *Decoder) decodeObject(dec Decodee) (err error) {
	var (
		// Byte currently being inspected
		b byte
		// State of our state machine
		state uint8
		// Value helper
	//	val Value
	)

	//	val.d = d

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		switch state {
		case osStart:
			if isWhitespace(b) {
				continue
			}

			if b != charDoubleQuote {
				err = ErrInvalidChar
				goto END
			}

			state = osKey

		case osKey:
			if b == charDoubleQuote {
				state = osPreSeparator
				continue
			}

			d.kb.WriteByte(b)

		case osPreSeparator:
			if isWhitespace(b) {
				continue
			}

			if b != charColon {
				err = ErrInvalidChar
				goto END
			}

			state = osValue

		case osValue:
			if d.v.vt, err = d.appendValue(b); err != nil {
				goto END
			}

			if err = dec.UnmarshalJsoon(unsafeString(d.kb.Bytes()), &d.v); err != nil {
				goto END
			}

			//val.vt = valNil
			d.v.vt = valNil
			d.kb.Reset()
			d.vb.Reset()
			state = osPostValue

		case osPostValue:
			if isWhitespace(b) {
				state = osEnd
			} else if b == charComma {
				state = osStart
			} else if b == charCloseCurly {
				state = osEnd
				goto END
			} else {
				err = ErrInvalidChar
				goto END
			}

		case osEnd:
			if isWhitespace(b) {
				continue
			}

			if b != charCloseCurly {
				err = ErrInvalidChar
			}

			goto END
		}
	}

END:

	if state != osEnd {
		return ErrUnexpectedEnd
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func (d *Decoder) decodeArray(dec ArrayDecodee) (err error) {
	var (
		// Byte currently being inspected
		b byte
		// State of our state machine
		state uint8
		// Value helper
	//	val Value
	)

	//val.d = d

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		switch state {
		case asStart:
			if isWhitespace(b) {
				continue
			}

			if d.v.vt, err = d.appendValue(b); err != nil {
				return
			}

			if err = dec.UnmarshalJsoon(&d.v); err != nil {
				return
			}

			//val.vt = valNil
			d.v.vt = valNil
			d.vb.Reset()
			state = asPostValue

		case asPostValue:
			if isWhitespace(b) {
				state = asEnd
			} else if b == charComma {
				state = asStart
			} else if b == charCloseBracket {
				state = asEnd
				goto END
			} else {
				err = ErrInvalidChar
				goto END
			}

		case asEnd:
			if isWhitespace(b) {
				continue
			} else if b == charCloseBracket {
				goto END
			} else {
				err = ErrInvalidChar
				goto END
			}
		}
	}

END:
	if state != asEnd {
		return ErrUnexpectedEnd
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func (d *Decoder) appendValue(lead byte) (vt uint8, err error) {
	var b byte
	for b = lead; err == nil; b, err = d.r.ReadByte() {
		if isWhitespace(b) {
			continue
		}

		switch b {
		case charDoubleQuote:
			vt = valString
			err = d.appendString()

		case charLowerT:
			vt = valBool
			err = d.appendTrue()

		case charLowerF:
			vt = valBool
			err = d.appendFalse()

		case charOpenCurly:
			vt = valObject
			return

		case charOpenBracket:
			vt = valArray
			return

		case charLowerN:
			vt = valNil
			err = d.readNull()
			return

		default:
			// TODO: Figure out a cleaner way to perform this check
			if isNumber(b) || b == charHyphen {
				vt = valNumber
				err = d.appendNumber(b)
			} else {
				err = ErrInvalidChar
			}
		}

		return
	}

	return
}

func (d *Decoder) appendString() (err error) {
	var (
		b       byte
		escaped bool
	)

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if b == charDoubleQuote && !escaped {
			return
		}

		if !escaped && b == charBackslash {
			escaped = true
			continue
		}

		d.vb.WriteByte(b)
		escaped = false
	}

	return ErrUnexpectedEnd
}

func (d *Decoder) appendNumber(lead byte) (err error) {
	var (
		b   byte
		cnt int
	)

	for b = lead; err == nil; b, err = d.r.ReadByte() {
		cnt++
		if isNumber(b) {
			d.vb.WriteByte(b)
			continue
		}

		switch b {
		case charSpace, charNewline, charTab:
			return
		case charComma, charCloseCurly, charCloseBracket:
			// TODO: Figure out a way to remove this UnreadByte
			d.r.UnreadByte()
			return
		default:
			if cnt == 1 && b == charHyphen {
				d.vb.WriteByte(b)
				continue
			}
			// Invalid character found, expected a number or a number-ending character
			return ErrInvalidChar
		}
	}

	// If we made it through the loop without finding the end to the number, we ended too early
	return ErrUnexpectedEnd
}

func (d *Decoder) appendTrue() (err error) {
	var b byte
	for i := 1; i < 4; i++ {
		if b, err = d.r.ReadByte(); err != nil {
			return
		} else if b != trueBytes[i] {
			return ErrInvalidChar
		}
	}

	d.vb.WriteBool(true)
	return
}

func (d *Decoder) appendFalse() (err error) {
	var b byte
	for i := 1; i < 5; i++ {
		if b, err = d.r.ReadByte(); err != nil {
			return
		} else if b != falseBytes[i] {
			return ErrInvalidChar
		}
	}

	d.vb.WriteBool(false)
	return
}

func (d *Decoder) readNull() (err error) {
	var b byte
	for i := 1; i < 4; i++ {
		if b, err = d.r.ReadByte(); err != nil {
			return
		} else if b != nullBytes[i] {
			return ErrInvalidChar
		}
	}

	return
}
