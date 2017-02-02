package parseFloat

import (
	"strconv"
	"testing"
)

var val float64
var testVal1 = []byte("123")
var testVal2 = []byte("-3.14159265e+23")

func TestBasic(t *testing.T) {
	if DecodeNumber(testVal1) != 123 {
		t.Fatal("testVal1 decoded incorrectly")
	}
	if DecodeNumber(testVal2) != -3.14159265e+23 {
		t.Fatal("testVal2 decoded incorrectly")
	}
}

func BenchmarkNumber_Val1(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		val = DecodeNumber(testVal1)
	}
}

func BenchmarkStdlibFloat_Val1(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		val, _ = strconv.ParseFloat(string(testVal1), 64)
	}
}

func BenchmarkNumber_Val2(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		val = DecodeNumber(testVal2)
	}
}

func BenchmarkStdlibFloat_Val2(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		val, _ = strconv.ParseFloat(string(testVal2), 64)
	}
}
