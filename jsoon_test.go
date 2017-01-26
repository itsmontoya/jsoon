package jsoon

import (
	"bytes"
	"encoding/json"
	"testing"
)

const (
	testStr = `{"name":"Test Name","greeting":"Hello world!","age":"32","activeUser":"true"}`
)

func TestMarshal(t *testing.T) {
	ts := newTestStruct()
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	enc.Encode(&ts)

	if str := buf.String(); str != testStr {
		t.Fatalf("invalid result\nExpected: %s\nReturned: %s\n", testStr, str)
	}
}

func BenchmarkJsoonMarshal(b *testing.B) {
	ts := newTestStruct()
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	enc := NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		enc.Encode(&ts)
		buf.Reset()
	}

	b.ReportAllocs()
}

func BenchmarkStdlibMarshal(b *testing.B) {
	ts := newTestStruct()
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	enc := json.NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		enc.Encode(&ts)
		buf.Reset()
	}

	b.ReportAllocs()
}

func newTestStruct() (ts testStruct) {
	ts.Name = "Test Name"
	ts.Greeting = "Hello world!"
	ts.Age = 32
	ts.ActiveUser = true
	return
}

type testStruct struct {
	Name       string
	Greeting   string
	Age        float64
	ActiveUser bool
}

func (t *testStruct) MarshalJsoon(enc *Encoder) (err error) {
	enc.String("name", t.Name)
	enc.String("greeting", t.Greeting)
	enc.Number("age", t.Age)
	enc.Bool("activeUser", t.ActiveUser)
	return
}
