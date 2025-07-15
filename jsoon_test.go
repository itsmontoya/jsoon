package jsoon

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/buger/jsonparser"
)

const (
	testStr = `{"name":"Test Name","greeting":"Hello \"world\"!","age":32,"activeUser":true,"additional":{"dateCreated":"2017-01-01","lastLogin":"2017-01-01"},"additionals":[{"dateCreated":"2017-01-01","lastLogin":"2017-01-01"},{"dateCreated":"2017-01-02","lastLogin":"2017-01-02"},{"dateCreated":"2017-01-03","lastLogin":"2017-01-03"}]}`

	testExpanded = `
{
	"name" : "Test Name",
	"greeting" : "Hello \"world\"!",
	"age" : 32,
	"activeUser" : true,
	"additional" : {
		"dateCreated" : "2017-01-01",
		"lastLogin" : "2017-01-01"
	},
	"additionals" : [
		{
			"dateCreated" : "2017-01-01",
			"lastLogin" : "2017-01-01"
		},
		{
			"dateCreated" : "2017-01-02",
			"lastLogin" : "2017-01-02"
		},
		{
			"dateCreated" : "2017-01-03",
			"lastLogin" : "2017-01-03"
		}
	]
}
`
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
	cts := newTestStruct()
	rdr := bytes.NewReader([]byte(testStr))
	dec := NewDecoder(rdr)

	// Test a normal decode process
	if err := dec.Decode(&ts); err != nil {
		t.Fatal(err)
	}

	// Compare values
	if !ts.Equals(&cts) {
		t.Fatalf("invalid value, expected <%+v> and received <%+v>", ts, cts)
	}

	// Test decoding again from the same reader
	rdr.Seek(0, 0)
	if err := dec.Decode(&ts); err != nil {
		t.Fatal(err)
	}

	// Compare values
	if !ts.Equals(&cts) {
		t.Fatal("invalid value")
	}

	// Test with the expanded object
	rdr = bytes.NewReader([]byte(testExpanded))
	dec = NewDecoder(rdr)
	if err := dec.Decode(&ts); err != nil {
		t.Fatal(err)
	}

	if !ts.Equals(&cts) {
		t.Fatal("invalid value")
	}
}

func TestStripeUnmarshal(t *testing.T) {
	var (
		sr  StripeChargeResponse
		f   *os.File
		err error
	)

	if f, err = os.Open("./testing/stripe.json"); err != nil {
		t.Fatal(err)
	}

	if err = NewDecoder(f).Decode(&sr); err != nil {
		t.Fatal(err)
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
	buf := bytes.NewReader([]byte(testStr))
	dec := NewDecoder(buf)

	for i := 0; i < b.N; i++ {
		dec.Decode(&ts)
		buf.Seek(0, 0)
	}

	b.ReportAllocs()
}

func BenchmarkJsoonUnmarshalPara(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var ts testStruct
		buf := bytes.NewReader([]byte(testStr))
		dec := NewDecoder(buf)

		for p.Next() {
			dec.Decode(&ts)
			buf.Seek(0, 0)
		}
	})

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
	buf := bytes.NewReader([]byte(testStr))
	dec := json.NewDecoder(buf)

	for i := 0; i < b.N; i++ {
		dec.Decode(&ts)
		buf.Seek(0, 0)
	}

	b.ReportAllocs()
}

func BenchmarkJsonParserUnmarshal(b *testing.B) {
	var ts testStruct
	buf := bytes.NewReader([]byte(testStr))

	for i := 0; i < b.N; i++ {
		// Need to read all to byteslice to simulate receiving a reader and needing to read all bytes before parsing
		data, err := io.ReadAll(buf)
		if err != nil {
			b.Fatal(err)
		}

		if ts.Name, err = jsonparser.GetString(data, "name"); err != nil {
			b.Fatal(err)
		}

		if ts.Greeting, err = jsonparser.GetString(data, "greeting"); err != nil {
			b.Fatal(err)
		}

		if ts.Age, err = jsonparser.GetFloat(data, "age"); err != nil {
			b.Fatal(err)
		}

		if ts.ActiveUser, err = jsonparser.GetBoolean(data, "activeUser"); err != nil {
			b.Fatal(err)
		}

		ts.Additional = &testSimpleStruct{}
		if ts.Additional.DateCreated, err = jsonparser.GetString(data, "additional", "dateCreated"); err != nil {
			b.Fatal(err)
		}

		if ts.Additional.LastLogin, err = jsonparser.GetString(data, "additional", "lastLogin"); err != nil {
			b.Fatal(err)
		}

		ts.Additionals = make(testSimpleStructSlice, 0, 3)
		jsonparser.ArrayEach(data, func(bs []byte, vt jsonparser.ValueType, offset int, err error) {
			tss := &testSimpleStruct{}
			if tss.DateCreated, err = jsonparser.GetString(bs, "dateCreated"); err != nil {
				b.Fatal(err)
			}

			if tss.LastLogin, err = jsonparser.GetString(bs, "lastLogin"); err != nil {
				b.Fatal(err)
			}

			ts.Additionals = append(ts.Additionals, tss)
		}, "additionals")

		buf.Seek(0, 0)
	}

	b.ReportAllocs()

}

func BenchmarkEasyJSONUnmarshal(b *testing.B) {
	var ts testStruct
	buf := bytes.NewReader([]byte(testStr))

	for i := 0; i < b.N; i++ {
		// Need to read all to byteslice to simulate receiving a reader and needing to read all bytes before parsing
		data, err := io.ReadAll(buf)
		if err != nil {
			b.Fatal(err)
		}

		if ts.Name, err = jsonparser.GetString(data, "name"); err != nil {
			b.Fatal(err)
		}

		if ts.Greeting, err = jsonparser.GetString(data, "greeting"); err != nil {
			b.Fatal(err)
		}

		if ts.Age, err = jsonparser.GetFloat(data, "age"); err != nil {
			b.Fatal(err)
		}

		if ts.ActiveUser, err = jsonparser.GetBoolean(data, "activeUser"); err != nil {
			b.Fatal(err)
		}

		ts.Additional = &testSimpleStruct{}
		if ts.Additional.DateCreated, err = jsonparser.GetString(data, "additional", "dateCreated"); err != nil {
			b.Fatal(err)
		}

		if ts.Additional.LastLogin, err = jsonparser.GetString(data, "additional", "lastLogin"); err != nil {
			b.Fatal(err)
		}

		ts.Additionals = make(testSimpleStructSlice, 0, 3)
		jsonparser.ArrayEach(data, func(bs []byte, vt jsonparser.ValueType, offset int, err error) {
			tss := &testSimpleStruct{}
			if tss.DateCreated, err = jsonparser.GetString(bs, "dateCreated"); err != nil {
				b.Fatal(err)
			}

			if tss.LastLogin, err = jsonparser.GetString(bs, "lastLogin"); err != nil {
				b.Fatal(err)
			}

			ts.Additionals = append(ts.Additionals, tss)
		}, "additionals")

		buf.Seek(0, 0)
	}

	b.ReportAllocs()

}

func newTestStruct() (ts testStruct) {
	ts.Name = "Test Name"
	ts.Greeting = `Hello "world"!`
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
	// We know our name will never have quotes, so we can utilize unsafe string
	enc.UnsafeString("name", t.Name)
	// Our greeting might have quotes, so we will ensure it's escaped just in case
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
		t.Additionals = make(testSimpleStructSlice, 0, 3)
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
	// We know our dateCreated and lastLogin will never have quotes, so we can utilize unsafe string
	enc.UnsafeString("dateCreated", t.DateCreated)
	enc.UnsafeString("lastLogin", t.LastLogin)
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
