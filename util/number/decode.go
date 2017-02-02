package parseFloat

import "math"

// DecodeNumber decodes a sequence of bytes representing a JSON Number into
// float64. Ported from cJSON.
//
// For reference see http://json.org/number.gif
func DecodeNumber(num []byte) (n float64) {
	var (
		sign         float64 = 1
		scale        int
		subscale     int
		signsubscale int = 1
	)

	func() {
		var i int

		// has sign?
		if num[i] == '-' {
			sign = -1
			i++
			if i >= len(num) {
				return
			}
		}
		// is zero
		if num[i] == '0' {
			i++
			if i >= len(num) {
				return
			}
		}
		// number?
		if num[i] >= '1' && num[i] <= '9' {
			for num[i] >= '0' && num[i] <= '9' {
				n = n*10 + float64(num[i]-'0')
				i++
				if i >= len(num) {
					return
				}
			}
		}
		// fractional part?
		if num[i] == '.' && (num[i+1] >= '0' && num[i+1] <= '9') {
			i++
			if i >= len(num) {
				return
			}
			for num[i] >= '0' && num[i] <= '9' {
				n = n*10 + float64(num[i]-'0')
				i++
				if i >= len(num) {
					return
				}
				scale--
			}
		}
		// exponent?
		if num[i] == 'e' || num[i] == 'E' {
			i++
			if i >= len(num) {
				return
			}
			// with sign?
			if num[i] == '+' {
				i++
				if i >= len(num) {
					return
				}
			} else if num[i] == '-' {
				signsubscale = -1
				i++
				if i >= len(num) {
					return
				}
			}
			// number?
			for num[i] >= '0' && num[i] <= '9' {
				subscale = subscale*10 + int(num[i]-'0')
				i++
				if i >= len(num) {
					return
				}
			}
		}
	}()

	// number = +/- number.fraction * 10^+/- exponent
	n = sign * n * math.Pow10(scale+subscale*signsubscale)

	return
}
