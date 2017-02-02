### Benchmarks

```
$ go test -benchtime=10s -bench=.
BenchmarkNumber_Val1-4          1000000000          18.3 ns/op         0 B/op          0 allocs/op
BenchmarkStdlibFloat_Val1-4     200000000           57.4 ns/op         3 B/op          1 allocs/op
BenchmarkNumber_Val2-4          300000000           43.1 ns/op         0 B/op          0 allocs/op
BenchmarkStdlibFloat_Val2-4     200000000           94.3 ns/op        16 B/op          1 allocs/op
PASS
ok      github.com/itsmontoya/jsoon    84.174s
```
