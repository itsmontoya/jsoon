package jsoon

import "io"
import "bufio"

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
}

func (d *Decoder) decodeObject(dec Decodee) (err error) {
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

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
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
			d.r.UnreadByte()
			if val, vt, err = d.appendValue(val); err != nil {
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
		// Value type
		vt  uint8
		val []byte
	)

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
			if val, vt, err = d.appendValue(val); err != nil {
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

	if err == io.EOF {
		err = nil
	}

	return
}

// Decode will decode
func (d *Decoder) Decode(value interface{}) (err error) {
	var b byte

	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if b == charSpace {
			continue
		}

		switch b {
		case charSpace:
		case charOpenCurly:
			dec, ok := value.(Decodee)
			if !ok {
				err = ErrInvalidChar
				return
			}

			d.r.UnreadByte()
			return d.decodeObject(dec)

		case charOpenBracket:
			dec, ok := value.(ArrayDecodee)
			if !ok {
				err = ErrInvalidChar
				return
			}

			d.r.UnreadByte()
			return d.decodeArray(dec)

		default:
			return ErrInvalidChar
		}

		err = ErrInvalidChar
		return
	}

	return
}

func (d *Decoder) appendValue(val []byte) (out []byte, vt uint8, err error) {
	var b byte
	for b, err = d.r.ReadByte(); err == nil; b, err = d.r.ReadByte() {
		if b == charSpace {
			continue
		}

		switch b {
		case charSpace:
			continue

		case charDoubleQuote:
			vt = valString
			out, err = d.appendString(val)

		case charLowerT:
			vt = valBool
			d.r.UnreadByte()
			out, err = d.appendTrue(val)

		case charLowerF:
			vt = valBool
			d.r.UnreadByte()
			out, err = d.appendFalse(val)

		case charOpenCurly:
			vt = valObject
			d.r.UnreadByte()
			out, err = d.appendObject(val)

		case charOpenBracket:
			vt = valArray
			d.r.UnreadByte()
			out, err = d.appendArray(val)

		default:
			// TODO: Figure out a cleaner way to perform this check
			if isNumber(b) {
				vt = valNumber
				d.r.UnreadByte()
				out, err = d.appendNumber(val)
			} else {
				panic("Unsupported type!")
			}
		}

		return
	}

	return
}
