package jsoon

import (
	"bufio"
	"io"
)

const (
	stateStart uint8 = iota
	stateObjStart
	stateKeyStart
	stateKey
	stateKeyEnd
	stateKVSeparator
	stateVal
	stateValEnd
	stateObjWaitingEnd
	stateObjEnd
	stateEnd
)

const (
	dsStart uint8 = iota
	dsObjectStart
	dsArrayStart
	dsEnd
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
	charOpenCurly    = '{'
	charCloseCurly   = '}'
	charOpenBracket  = '['
	charCloseBracket = ']'
	charColon        = ':'
	charComma        = ','
	charZero         = '0'
	charNine         = '9'
	charLowerA       = 'a'
	charLowerZ       = 'z'
	charUpperA       = 'A'
	charUpperZ       = 'Z'
	charLowerT       = 't'
	charLowerF       = 'f'
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

			d.r.UnreadByte()
			err = d.decodeObject(dec)
			goto END

		case charOpenBracket:
			dec, ok := value.(ArrayDecodee)
			if !ok {
				err = ErrInvalidChar
				goto END
			}

			d.r.UnreadByte()
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
		val Value
	)

	val.d = d

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		switch state {
		case osStart:
			if isWhitespace(b) {
				continue
			}

			if b == charDoubleQuote {
				state = osKey
			}

		case osKey:
			if b == charDoubleQuote {
				state = osPreSeparator
				continue
			}

			d.kb.WriteByte(b)

		case osPreSeparator:
			if b == charSpace {
				continue
			}

			if b == charColon {
				state = osValue
			} else {
				err = ErrInvalidChar
				return
			}

		case osValue:
			d.r.UnreadByte()
			if val.vt, err = d.appendValue(); err != nil {
				return
			}

			if err = dec.UnmarshalJsoon(string(d.kb.Bytes()), &val); err != nil {
				return
			}

			val.vt = valNil
			d.kb.Reset()
			d.vb.Reset()
			state = osPostValue

		case osPostValue:
			switch b {
			case charComma:
				state = osStart

			case charCloseCurly:
				state = osEnd
				goto END

			case charSpace, charNewline, charTab:
				state = osEnd

			default:
				return ErrInvalidChar
			}

		case osEnd:
			switch b {
			case charSpace, charNewline, charTab:
			case charCloseCurly:
				goto END

			default:
				return ErrInvalidChar
			}
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
		val Value
	)

	val.d = d

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		switch state {
		case asStart:
			if b == charOpenBracket {
				state = asValue
				continue
			}

			return ErrInvalidChar

		case asValue:
			d.r.UnreadByte()
			if val.vt, err = d.appendValue(); err != nil {
				return
			}

			if err = dec.UnmarshalJsoon(&val); err != nil {
				return
			}

			val.vt = valNil
			d.vb.Reset()
			state = asPostValue

		case asPostValue:
			switch b {
			case charComma:
				state = asValue

			case charCloseBracket:
				state = asEnd
				goto END

			case charSpace, charNewline, charTab:
				state = asEnd

			default:
				return ErrInvalidChar
			}

		case asEnd:
			switch b {
			case charSpace, charNewline, charTab:
			case charCloseBracket:
				goto END

			default:
				return ErrInvalidChar
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

func (d *Decoder) appendValue() (vt uint8, err error) {
	var b byte
	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		switch b {
		case charSpace, charNewline, charTab:
			continue

		case charDoubleQuote:
			vt = valString
			err = d.appendString()

		case charLowerT:
			vt = valBool
			d.r.UnreadByte()
			err = d.appendTrue()

		case charLowerF:
			vt = valBool
			d.r.UnreadByte()
			err = d.appendFalse()

		case charOpenCurly:
			vt = valObject
			d.r.UnreadByte()
			return

		case charOpenBracket:
			vt = valArray
			d.r.UnreadByte()
			return

		default:
			// TODO: Figure out a cleaner way to perform this check
			if isNumber(b) {
				vt = valNumber
				d.r.UnreadByte()
				err = d.appendNumber()
			} else {
				panic("Unsupported type!")
			}
		}

		return
	}

	return
}

func (d *Decoder) appendString() (err error) {
	var b byte
	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if b == charDoubleQuote {
			return
		}

		d.vb.WriteByte(b)
	}

	return ErrUnexpectedEnd
}

func (d *Decoder) appendNumber() (err error) {
	var b byte
	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if isNumber(b) {
			d.vb.WriteByte(b)
			continue
		}

		switch b {
		case charSpace, charNewline, charTab, charComma, charCloseCurly:
			d.r.UnreadByte()
			return
		default:
			// Invalid character found, expected a number or a number-ending character
			return ErrInvalidChar
		}
	}

	// If we made it through the loop without finding the end to the number, we ended too early
	return ErrUnexpectedEnd
}

func (d *Decoder) appendTrue() (err error) {
	var b byte
	for i := 0; i < 4; i++ {
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
	for i := 0; i < 5; i++ {
		if b, err = d.r.ReadByte(); err != nil {
			return
		} else if b != falseBytes[i] {
			return ErrInvalidChar
		}
	}

	d.vb.WriteBool(false)
	return
}

func (d *Decoder) appendObject() (err error) {
	var (
		b     byte
		depth int
	)

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		d.vb.WriteByte(b)
		switch b {
		case charCloseCurly:
			if depth--; depth == 0 {
				return
			}

		case charOpenCurly:
			depth++
		}
	}

	return ErrUnexpectedEnd
}

func (d *Decoder) appendArray() (err error) {
	var (
		b     byte
		depth int
	)

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		d.vb.WriteByte(b)
		switch b {
		case charCloseBracket:
			if depth--; depth == 0 {
				return
			}
		case charOpenBracket:
			depth++
		}
	}

	return ErrUnexpectedEnd
}
