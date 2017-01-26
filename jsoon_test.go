package jsoon

import (
	"bytes"
	"encoding/json"
	"testing"
)

const (
	testStr = `{"name":"Test Name","greeting":"Hello world!","age":32,"activeUser":true,"additional":{"dateCreated":"2017-01-01","lastLogin":"2017-01-01"}}`
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
	ts.Additional.DateCreated = "2017-01-01"
	ts.Additional.LastLogin = "2017-01-01"
	return
}

type testStruct struct {
	Name       string
	Greeting   string
	Age        float64
	ActiveUser bool

	Additional testSimpleStruct
}

func (t *testStruct) MarshalJsoon(enc *Encoder) (err error) {
	enc.String("name", t.Name)
	enc.String("greeting", t.Greeting)
	enc.Number("age", t.Age)
	enc.Bool("activeUser", t.ActiveUser)
	enc.Object("additional", &t.Additional)
	return
}

type testSimpleStruct struct {
	DateCreated string
	LastLogin   string
}

func (t *testSimpleStruct) MarshalJsoon(enc *Encoder) (err error) {
	enc.String("dateCreated", t.DateCreated)
	enc.String("lastLogin", t.LastLogin)
	return
}
