## Jsoon 
Jsoon is a fast and simple json encoding/decoding library. Custom Marshal/Unmarshal funcs are utilized to avoid leveraging reflection.

*At the moment, only encoding is implemented*

## Benchmarks
```
BenchmarkJsoonMarshal-4          5000000               365 ns/op               0 B/op          0 allocs/op
BenchmarkStdlibMarshal-4         1000000              1337 ns/op               8 B/op          1 allocs/op
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