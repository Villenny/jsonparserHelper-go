[![GitHub issues](https://img.shields.io/github/issues/Villenny/jsonparserHelper-go)](https://github.com/Villenny/jsonparserHelper-go/issues)
[![GitHub forks](https://img.shields.io/github/forks/Villenny/jsonparserHelper-go)](https://github.com/Villenny/jsonparserHelper-go/network)
[![GitHub stars](https://img.shields.io/github/stars/Villenny/jsonparserHelper-go)](https://github.com/Villenny/jsonparserHelper-go/stargazers)
[![GitHub license](https://img.shields.io/github/license/Villenny/jsonparserHelper-go)](https://github.com/Villenny/jsonparserHelper-go/blob/master/LICENSE)
![Go](https://github.com/Villenny/jsonparserHelper-go/workflows/Go/badge.svg?branch=master)
![Codecov branch](https://img.shields.io/codecov/c/github/villenny/jsonparserHelper-go/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Villenny/jsonparserHelper-go)](https://goreportcard.com/report/github.com/Villenny/jsonparserHelper-go)
[![Documentation](https://godoc.org/github.com/Villenny/jsonparserHelper-go?status.svg)](http://godoc.org/github.com/Villenny/jsonparserHelper-go)

# jsonparserHelper-go
- zero allocation json parsing
- convenience helper for serializing with buger/jsonpaser to a struct


## Install

```
go get -u github.com/Villenny/jsonparserHelper-go
```

## Notable members:
`MakeParser`,
`RunParser`,

The expected use case:
- declare your parser somewhere
```
	import (
		"github.com/buger/jsonparser"
		"github.com/villenny/jsonparserHelper-go"
	)

	type Fubar struct {
		Id       string
		BidId    string
		BidImpId string
		BidPrice float64
		BidAdm   string
		Cur      string
		Foo      string
	}

	var fubarParser = func() []Parser {
		parser := MakeParser(make([]Parser, 32)[:0]).
			Add(UnsafeStringValue(unsafe.Offsetof(Fubar{}.Id)), "id").
			Add(UnsafeStringValue(unsafe.Offsetof(Fubar{}.BidId)), "seatbid", "[0]", "bid", "[0]", "id").
			Add(UnsafeStringValue(unsafe.Offsetof(Fubar{}.BidImpId)), "seatbid", "[0]", "bid", "[0]", "impid").
			Add(Float64Value(unsafe.Offsetof(Fubar{}.BidPrice)), "seatbid", "[0]", "bid", "[0]", "price").
			Add(UnsafeStringValue(unsafe.Offsetof(Fubar{}.BidAdm)), "seatbid", "[0]", "bid", "[0]", "adm").
			Add(UnsafeStringValue(unsafe.Offsetof(Fubar{}.Cur)), "cur").
			Add(UnsafeStringValue(unsafe.Offsetof(Fubar{}.Foo)), "foo")
		return parser
	}()
```

Use your parser. Yay zero allocations (almost, buger/jsonparser still has a few but I assume those will be fixed soon)
```
		fubar := Fubar{}
		err := RunParser(bodyBytes, fubarParser, unsafe.Pointer(&fubar))
```


## Benchmark

- Still have 3 allocs inside buger/jsonparser EachKey function, once those are fixed (I contributed a PR which is under review), will be zero

```
$ ./bench.sh
=== RUN   TestRunParser
=== RUN   TestRunParser/simple_test
--- PASS: TestRunParser (0.00s)
    --- PASS: TestRunParser/simple_test (0.00s)
goos: windows
goarch: amd64
pkg: github.com/villenny/jsonparserHelper-go
BenchmarkRunParser
BenchmarkRunParser-8      924062              3947 ns/op              96 B/op          3 allocs/op
BenchmarkEachKey
BenchmarkEachKey-8        973983              3854 ns/op              96 B/op          3 allocs/op
PASS
ok      github.com/villenny/jsonparserHelper-go 7.701s

```

## Contact

Ryan Haksi [ryan.haksi@gmail.com]

## License

Available under the BSD [License](/LICENSE). Or any license really, do what you like

