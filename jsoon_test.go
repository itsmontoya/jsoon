package jsoon

import (
	"bytes"
	"encoding/json"
	"testing"
)

const (
	testStr = `{"name":"Test Name","greeting":"Hello world!","age":32,"activeUser":true,"additional":{"dateCreated":"2017-01-01","lastLogin":"2017-01-01"},"additionals":[{"dateCreated":"2017-01-01","lastLogin":"2017-01-01"},{"dateCreated":"2017-01-02","lastLogin":"2017-01-02"},{"dateCreated":"2017-01-03","lastLogin":"2017-01-03"}]}`
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
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	enc := NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		enc.Encode(&ts)
		buf.Reset()
	}

	b.ReportAllocs()
}

func BenchmarkStdlibMarshal(b *testing.B) {
	ts := newTestStruct()
	buf := bytes.NewBuffer(make([]byte, 0, 512))
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
	ts.Additional = &testSimpleStruct{
		DateCreated: "2017-01-01",
		LastLogin:   "2017-01-01",
	}

	ts.Additionals = append(ts.Additionals, &testSimpleStruct{"2017-01-01", "2017-01-01"})
	ts.Additionals = append(ts.Additionals, &testSimpleStruct{"2017-01-02", "2017-01-02"})
	ts.Additionals = append(ts.Additionals, &testSimpleStruct{"2017-01-03", "2017-01-03"})
	return
}

type testStruct struct {
	Name       string
	Greeting   string
	Age        float64
	ActiveUser bool

	Additional  *testSimpleStruct
	Additionals testSimpleStructSlice
}

func (t *testStruct) MarshalJsoon(enc *Encoder) (err error) {
	enc.String("name", t.Name)
	enc.String("greeting", t.Greeting)
	enc.Number("age", t.Age)
	enc.Bool("activeUser", t.ActiveUser)
	enc.Object("additional", t.Additional)
	enc.Array("additionals", t.Additionals)
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

type testSimpleStructSlice []*testSimpleStruct

func (t testSimpleStructSlice) MarshalJsoon(a *ArrayEncoder) (err error) {
	for _, v := range t {
		a.Object(v)
	}

	return
}
