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

func TestUnmarshal(t *testing.T) {
	var ts testStruct
	dec := NewDecoder(bytes.NewReader([]byte(testStr)))
	if err := dec.Decode(&ts); err != nil {
		t.Fatal(err)
	}

	cts := newTestStruct()
	if !ts.Equals(&cts) {
		t.Fatal("invalid value")
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

func BenchmarkJsoonUnmarshal(b *testing.B) {
	var ts testStruct

	for i := 0; i < b.N; i++ {
		NewDecoder(bytes.NewReader([]byte(testStr))).Decode(&ts)
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

func BenchmarkStdlibUnmarshal(b *testing.B) {
	var ts testStruct

	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(testStr), &ts)
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

func (t *testStruct) Equals(b *testStruct) bool {
	if t.Name != b.Name {
		return false
	}

	if t.Greeting != b.Greeting {
		return false
	}

	if t.Age != b.Age {
		return false
	}

	if t.ActiveUser != b.ActiveUser {
		return false
	}

	if !t.Additional.Equals(b.Additional) {
		return false
	}

	if len(t.Additionals) != len(b.Additionals) {
		return false
	}

	for i, ts := range t.Additionals {
		if !ts.Equals(b.Additionals[i]) {
			return false
		}
	}

	return true
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

func (t *testStruct) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "name":
		if t.Name, err = val.String(); err != nil {
			return
		}

	case "greeting":
		if t.Greeting, err = val.String(); err != nil {
			return
		}

	case "age":
		if t.Age, err = val.Number(); err != nil {
			return
		}

	case "activeUser":
		if t.ActiveUser, err = val.Bool(); err != nil {
			return
		}

	case "additional":
		t.Additional = &testSimpleStruct{}
		if err = val.Object(t.Additional); err != nil {
			return
		}

	case "additionals":
		t.Additionals = make(testSimpleStructSlice, 0)
		if err = val.Array(&t.Additionals); err != nil {
			return
		}
	}

	return
}

type testSimpleStruct struct {
	DateCreated string
	LastLogin   string
}

func (t *testSimpleStruct) Equals(b *testSimpleStruct) bool {
	if t.DateCreated != b.DateCreated {
		return false
	}

	if t.LastLogin != b.LastLogin {
		return false
	}

	return true
}

func (t *testSimpleStruct) MarshalJsoon(enc *Encoder) (err error) {
	enc.String("dateCreated", t.DateCreated)
	enc.String("lastLogin", t.LastLogin)
	return
}

func (t *testSimpleStruct) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "dateCreated":
		if t.DateCreated, err = val.String(); err != nil {
			return
		}
	case "lastLogin":
		if t.LastLogin, err = val.String(); err != nil {
			return
		}
	}

	return
}

type testSimpleStructSlice []*testSimpleStruct

func (t testSimpleStructSlice) MarshalJsoon(a *ArrayEncoder) (err error) {
	for _, v := range t {
		a.Object(v)
	}

	return
}

func (t *testSimpleStructSlice) UnmarshalJsoon(val *Value) (err error) {
	var ts testSimpleStruct
	if err = val.Object(&ts); err != nil {
		return
	}

	*t = append(*t, &ts)
	return
}
