package jsoon

import "io"

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
func NewDecoder(s []byte) *Decoder {
	return &Decoder{
		s: s,
	}
}

// Decoder handles decoding
type Decoder struct {
	r io.Reader
	s []byte
}

func (d *Decoder) decodeObject(dec Decodee, n int) (err error) {
	var (
		// Byte currently being inspected
		b byte
		// State of our state machine
		state uint8
		// Value type
		vt uint8

		key []byte
		val []byte
	)

	for i := n; i < len(d.s); i++ {
		b = d.s[i]

		switch state {
		case osStart:
			if b == charSpace {
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

			key = append(key, b)

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
			var n int
			if val, vt, n = d.appendValue(val, d.s[i:]); n == -1 {
				err = ErrInvalidChar
				return
			}

			if err = dec.UnmarshalJsoon(string(key), &Value{
				vt: vt,
				d:  val,
			}); err != nil {
				return
			}

			vt = valNil
			key = key[:0]
			val = val[:0]
			state = osPostValue
			i += n

		case osPostValue:
			switch b {
			case charComma:
				state = osStart

			case charCloseCurly:
				state = osEnd
				break

			case charSpace:
				state = osEnd

			default:
				return ErrInvalidChar
			}

		case osEnd:
			switch b {
			case charSpace:
			case charCloseCurly:
				break

			default:
				return ErrInvalidChar
			}
		}
	}

	if state != osEnd {
		return ErrUnexpectedEnd
	}

	return
}

func (d *Decoder) decodeArray(dec ArrayDecodee, n int) (err error) {
	var (
		// Byte currently being inspected
		b byte
		// State of our state machine
		state uint8
		// Value type
		vt  uint8
		val []byte
	)

	for i := n; i < len(d.s); i++ {
		b = d.s[i]

		switch state {
		case asStart:
			if b == charOpenBracket {
				state = asValue
				continue
			}

			return ErrInvalidChar

		case asValue:
			var n int
			if val, vt, n = d.appendValue(val, d.s[i:]); n == -1 {
				err = ErrInvalidChar
				return
			}

			if err = dec.UnmarshalJsoon(&Value{
				vt: vt,
				d:  val,
			}); err != nil {
				return
			}

			vt = valNil
			val = val[:0]
			state = asPostValue
			i += n

		case asPostValue:
			switch b {
			case charComma:
				state = asValue

			case charCloseBracket:
				state = asEnd
				break

			case charSpace:
				state = asEnd

			default:
				return ErrInvalidChar
			}

		case asEnd:
			switch b {
			case charSpace:
			case charCloseCurly:
				break

			default:
				return ErrInvalidChar
			}
		}
	}

	if state != asEnd {
		return ErrUnexpectedEnd
	}

	return
}

// Decode will decode
func (d *Decoder) Decode(value interface{}) (err error) {
	var (
		// Byte currently being inspected
		b byte
		// State of our state machine
		state uint8
	)

	for i := 0; i < len(d.s); i++ {
		b = d.s[i]
		switch state {
		case dsStart:
			if b == charSpace {
				continue
			}

			if b == charOpenCurly {
				state = dsObjectStart
				i--
				continue
			}

			if b == charOpenBracket {
				state = dsArrayStart
				i--
				continue
			}

			err = ErrInvalidChar
			return

		case dsObjectStart:
			dec, ok := value.(Decodee)
			if !ok {
				err = ErrInvalidChar
				return
			}

			return d.decodeObject(dec, i)
		case dsArrayStart:
			dec, ok := value.(ArrayDecodee)
			if !ok {
				err = ErrInvalidChar
				return
			}

			return d.decodeArray(dec, i)
		}
	}

	return
}

func (d *Decoder) appendValue(val, s []byte) (out []byte, vt uint8, n int) {
	for i, b := range s {
		if b == charSpace {
			continue
		}

		switch b {
		case charSpace:
			continue

		case charDoubleQuote:
			vt = valString
			out, n = appendString(val, s[i:])

		case charLowerT, charLowerF:
			vt = valBool
			out, n = appendBool(val, s[i:])

		case charOpenCurly:
			vt = valObject
			out, n = appendObject(val, s[i:])

		case charOpenBracket:
			vt = valArray
			out, n = appendArray(val, s[i:])

		default:
			// TODO: Figure out a cleaner way to perform this check
			if isNumber(b) {
				vt = valNumber
				out, n = appendNumber(val, s[i:])
			} else {
				panic("Unsupported type!")
			}
		}

		return
	}

	return val, vt, -1
}

/*
Reading func, will approach later
// Decode will decode
func (d *Decoder) Decode(value Decodee) (err error) {
	var buf [32]byte
	var state uint8
	var key []byte
	//	var val []byte
	var vt uint8

	for n, err := d.r.Read(buf[:]); err == nil; n, err = d.r.Read(buf[:]) {
		s := buf[:n]

		for i := 0; i < len(s); i++ {

		}
	}

	return
}
*/
