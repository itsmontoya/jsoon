package jsoon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"os"
	"testing"
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
		t.Fatal("invalid value")
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
		data, err := ioutil.ReadAll(buf)
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

type StripeChargeResponse struct {
	ID                 string         `json:"id"`
	Object             string         `json:"object"`
	Amount             int            `json:"amount"`
	AmountRefunded     int            `json:"amount_refunded"`
	BalanceTransaction string         `json:"balance_transaction"`
	Captured           bool           `json:"captured"`
	Created            int            `json:"created"`
	Currency           string         `json:"currency"`
	Customer           string         `json:"customer"`
	Livemode           bool           `json:"livemode"`
	Outcome            *StripeOutcome `json:"outcome"`
	Paid               bool           `json:"paid"`
	Refunded           bool           `json:"refunded"`
	Refunds            *StripeRefunds `json:"refunds"`
	Source             *StripeSource  `json:"source"`
	Status             string         `json:"status"`
}

func (s *StripeChargeResponse) MarshalJsoon(e *Encoder) (err error) {
	e.String("id", s.ID)
	e.String("object", s.Object)
	e.Number("amount", float64(s.Amount))
	e.Number("amount_refunded", float64(s.AmountRefunded))
	e.String("balance_transaction", s.BalanceTransaction)
	e.Bool("captured", s.Captured)
	e.Number("created", float64(s.Created))
	e.String("currency", s.Currency)
	e.String("customer", s.Customer)
	e.Bool("livemode", s.Livemode)
	e.Object("outcome", s.Outcome)
	e.Bool("paid", s.Paid)
	e.Bool("refunded", s.Refunded)
	e.Object("refunds", s.Refunds)
	e.Object("source", s.Source)
	e.String("status", s.Status)
	return
}

func (s *StripeChargeResponse) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "id":
		s.ID, err = val.String()
	case "object":
		s.Object, err = val.String()
	case "amount":
		var a float64
		a, err = val.Number()
		s.Amount = int(a)
	case "amount_refunded":
		var ar float64
		ar, err = val.Number()
		s.AmountRefunded = int(ar)
	case "balance_transaction":
		s.BalanceTransaction, err = val.String()
	case "captured":
		s.Captured, err = val.Bool()
	case "created":
		var cr float64
		cr, err = val.Number()
		s.Created = int(cr)
	case "currency":
		s.Currency, err = val.String()
	case "customer":
		s.Customer, err = val.String()
	case "livemode":
		s.Livemode, err = val.Bool()
	case "outcome":
		s.Outcome = &StripeOutcome{}
		err = val.Object(s.Outcome)
	case "paid":
		s.Paid, err = val.Bool()
	case "refunded":
		s.Refunded, err = val.Bool()
	case "refunds":
		s.Refunds = &StripeRefunds{}
		err = val.Object(s.Refunds)
	case "source":
		s.Source = &StripeSource{}
		err = val.Object(s.Source)
	case "status":
		s.Status, err = val.String()
	}

	return
}

func (s *StripeChargeResponse) Equals(os *StripeChargeResponse) (err error) {
	return
}

type StripeOutcome struct {
	NetworkStatus string `json:"network_status"`
	RiskLevel     string `json:"risk_level"`
	SellerMessage string `json:"seller_message"`
	Type          string `json:"type"`
}

func (s *StripeOutcome) MarshalJsoon(e *Encoder) (err error) {
	e.String("network_status", s.NetworkStatus)
	e.String("risk_level", s.RiskLevel)
	e.String("seller_message", s.SellerMessage)
	e.String("type", s.Type)
	return
}

func (s *StripeOutcome) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "network_status":
		s.NetworkStatus, err = val.String()
	case "risk_level":
		s.RiskLevel, err = val.String()
	case "seller_message":
		s.SellerMessage, err = val.String()
	case "type":
		s.Type, err = val.String()
	}

	return
}

func (s *StripeOutcome) Equals(os *StripeOutcome) (err error) {
	return
}

type StripeRefunds struct {
	Object     string `json:"object"`
	HasMore    bool   `json:"has_more"`
	TotalCount int    `json:"total_count"`
	URL        string `json:"url"`
}

func (s *StripeRefunds) MarshalJsoon(e *Encoder) (err error) {
	e.String("object", s.Object)
	e.Bool("has_more", s.HasMore)
	e.Number("total_count", float64(s.TotalCount))
	e.String("url", s.URL)
	return
}

func (s *StripeRefunds) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "object":
		s.Object, err = val.String()
	case "has_more":
		s.HasMore, err = val.Bool()
	case "total_count":
		var tc float64
		tc, err = val.Number()
		s.TotalCount = int(tc)
	case "url":
		s.URL, err = val.String()
	}

	return
}

func (s *StripeRefunds) Equals(os *StripeRefunds) (err error) {
	return
}

type StripeSource struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Brand    string `json:"brand"`
	Country  string `json:"country"`
	Customer string `json:"customer"`
	CVCCheck string `json:"cvc_check"`
	ExpMonth int    `json:"exp_month"`
	ExpYear  int    `json:"exp_year"`
	Funding  string `json:"funding"`
	Last4    string `json:"last4"`
}

func (s *StripeSource) MarshalJsoon(e *Encoder) (err error) {
	e.String("id", s.ID)
	e.String("object", s.Object)
	e.String("brand", s.Brand)
	e.String("country", s.Country)
	e.String("customer", s.Customer)
	e.String("cvc_check", s.CVCCheck)
	e.Number("exp_month", float64(s.ExpMonth))
	e.Number("exp_year", float64(s.ExpYear))
	e.String("funding", s.Funding)
	e.String("last4", s.Last4)
	return
}

func (s *StripeSource) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "id":
		s.ID, err = val.String()
	case "object":
		s.Object, err = val.String()
	case "brand":
		s.Brand, err = val.String()
	case "country":
		s.Country, err = val.String()
	case "customer":
		s.Customer, err = val.String()
	case "cvc_check":
		s.CVCCheck, err = val.String()
	case "exp_month":
		var em float64
		em, err = val.Number()
		s.ExpMonth = int(em)
	case "exp_year":
		var ey float64
		ey, err = val.Number()
		s.ExpYear = int(ey)
	case "funding":
		s.Funding, err = val.String()
	case "last4":
		s.Last4, err = val.String()
	}

	return
}

func (s *StripeSource) Equals(os *StripeSource) (err error) {
	if s.ID != os.ID {
		return fmt.Errorf("id's don't match: %s | %s")
	}

	if s.Object != os.Object {
		///		// return fmt.Errorf("objects's don't match: %s | %s")
	}

	if s.Brand != os.Brand {
		//	// return fmt.Errorf("brands's don't match: %s | %s")
	}

	if s.Country != os.Country {
		//	// return fmt.Errorf("countries's don't match: %s | %s")
	}

	if s.Customer != os.Customer {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	if s.CVCCheck != os.CVCCheck {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	if s.ExpMonth != os.ExpMonth {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	if s.ExpYear != os.ExpYear {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	if s.Last4 != os.Last4 {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	return
}
