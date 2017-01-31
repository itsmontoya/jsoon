## Jsoon [![GoDoc](https://godoc.org/github.com/itsmontoya/jsoon?status.svg)](https://godoc.org/github.com/itsmontoya/jsoon) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)

Jsoon is a fast and simple json encoding/decoding library. Custom Marshal/Unmarshal funcs are utilized to avoid leveraging reflection.

*At the moment, only encoding is implemented*

## Benchmarks
```
BenchmarkJsoonMarshal-4      1000000    1177 ns/op     40 B/op     2 allocs/op
BenchmarkJsoonUnmarshal-4     300000    6161 ns/op    632 B/op    38 allocs/op
BenchmarkStdlibMarshal-4      500000    2580 ns/op      8 B/op     1 allocs/op
BenchmarkStdlibUnmarshal-4    200000    9276 ns/op    160 B/op    11 allocs/op


```

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
2. Optimize the decoding process so that we incur less allocations
3. ???
4. Profit