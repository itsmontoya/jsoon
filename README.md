## Jsoon [![GoDoc](https://godoc.org/github.com/itsmontoya/jsoon?status.svg)](https://godoc.org/github.com/itsmontoya/jsoon) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)

Jsoon is a fast and simple json encoding/decoding library. Custom Marshal/Unmarshal funcs are utilized to avoid leveraging reflection.

*At the moment, only encoding is implemented*

## Benchmarks
```bash
## Go 1.7.4
# Jsoon
BenchmarkJsoonMarshal-4            2000000       691 ns/op        32 B/op       1 allocs/op
BenchmarkJsoonUnmarshal-4           500000      3554 ns/op       320 B/op      15 allocs/op
# Jsoon concurrent unmarshal
BenchmarkJsoonUnmarshalPara-4      2000000       941 ns/op       320 B/op      15 allocs/op
# Standard library
BenchmarkStdlibMarshal-4            500000      2637 ns/op         8 B/op       1 allocs/op
BenchmarkStdlibUnmarshal-4          200000      8952 ns/op       160 B/op      11 allocs/op
# github.com/buger/jsonparser
BenchmarkJsonParserUnmarshal-4      300000      4610 ns/op      2368 B/op      17 allocs/op

## Go 1.8 rc3
# Jsoon
BenchmarkJsoonMarshal-4            2000000       703 ns/op        32 B/op       1 allocs/op
BenchmarkJsoonUnmarshal-4           500000      3627 ns/op       320 B/op      15 allocs/op
# Jsoon concurrent unmarshal
BenchmarkJsoonUnmarshalPara-4      2000000       870 ns/op       320 B/op      15 allocs/op
# Standard library
BenchmarkStdlibMarshal-4           1000000      2321 ns/op         8 B/op       1 allocs/op
BenchmarkStdlibUnmarshal-4          200000      9361 ns/op       160 B/op      11 allocs/op
# github.com/buger/jsonparser
BenchmarkJsonParserUnmarshal-4      300000      4389 ns/op      2368 B/op      17 allocs/op

```

*Note: Even though the provided bench shows 1.8 as slower, on average 1.8 is faster AND more consistent*

## Usage 
```go
// See examples/webserver/webserver.go for source code
package main

import (
	"net/http"
	"sync"

	"github.com/itsmontoya/jsoon"
)

func main() {
	var s srv
	// Pre-set some values
	s.u.name = "Panda"
	s.u.age = 30

	err := http.ListenAndServe(":8080", &s)
	if err != nil {
		panic(err)
	}
}

type srv struct {
	u user
}

func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jsoon.NewEncoder(w).Encode(&s.u)
}

type user struct {
	mux sync.RWMutex

	name string
	age  int
}

func (u *user) MarshalJsoon(enc *jsoon.Encoder) (err error) {
	// We can lock from within the struct to ensure thread safety
	u.mux.RLock()
	enc.String("name", u.name)
	enc.Number("age", float64(u.age))
	u.mux.RUnlock()
	return
}
```

## To do
1. Add some more thorough testing, preferably using some JSON objects from common open APIs
