package jsonparserHelper

import (
	"testing"
	"unsafe"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

type TestResult struct {
	us  string
	i64 int64
	i   int
	f   float64
}

func TestRunParser(t *testing.T) {
	t.Run("simple test", func(t *testing.T) {
		parser := MakeParser(make([]Parser, 32)[:0]).
			Add(UnsafeStringValue(unsafe.Offsetof(TestResult{}.us)), "us").
			Add(Int64Value(unsafe.Offsetof(TestResult{}.i64)), "i64").
			Add(IntValue(unsafe.Offsetof(TestResult{}.i)), "i").
			Add(Float64Value(unsafe.Offsetof(TestResult{}.f)), "f")

		var results TestResult
		err := RunParser([]byte(`{
			"us": "some string",
			"i64": "32"
			"i": "64"
			"f": "3.14"
		}`), parser, unsafe.Pointer(&results))

		assert.Equal(t, "some string", results.us)
		assert.Equal(t, int64(32), results.i64)
		assert.Equal(t, int(64), results.i)
		assert.Equal(t, float64(3.14), results.f)
		assert.Nil(t, err)
	})
}

type Fubar struct {
	Id       string
	BidId    string
	BidImpId string
	BidPrice float64
	BidAdm   string
	Cur      string
	Foo      string
}

/*
	{"id"},
	{"seatbid", "[0]", "bid", "[0]", "id"},
	{"seatbid", "[0]", "bid", "[0]", "impid"},
	{"seatbid", "[0]", "bid", "[0]", "price"},
	{"seatbid", "[0]", "bid", "[0]", "adm"},
	{"cur"},
	{"foo"},
*/

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

/*
Deserialize with:
		fubar := Fubar{}
		err := RunParser(bodyBytes, fubarParser, unsafe.Pointer(&fubar))
*/

/*
BenchmarkRunParser-8
  308013	      3947 ns/op	      96 B/op	       3 allocs/op
PASS
*/
func BenchmarkRunParser(b *testing.B) {

	bodyBytes := StringAsBytes(BidResponse_Banner)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fubar := Fubar{Foo: "boo!"}
		err := RunParser(bodyBytes, fubarParser, unsafe.Pointer(&fubar))

		if "9f3701ef-1d8a-4b67-a601-e9a1a2fbcb7c" != fubar.Id ||
			"USD" != fubar.Cur ||
			"0" != fubar.BidId ||
			"1" != fubar.BidImpId ||
			0.35258 != fubar.BidPrice ||
			`some \"html\" code\nnext line` != fubar.BidAdm ||
			"boo!" != fubar.Foo ||
			err != nil {
			panic("ack")
		}

	}
}

/*
func BenchmarkParse(b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// for some reason go thinks this will escape if I define paths here
		// paths := BidResponse_Banner_Paths

		bodyBytes := StringAsBytes(BidResponse_Banner)

		results := make([]JsonparserResult, 7)[:]
		_ = Parse(bodyBytes, results, [][]string{
			{"id"},
			{"seatbid", "[0]", "bid", "[0]", "id"},
			{"seatbid", "[0]", "bid", "[0]", "impid"},
			{"seatbid", "[0]", "bid", "[0]", "price"},
			{"seatbid", "[0]", "bid", "[0]", "adm"},
			{"cur"},
			{"foo"},
		}...)
		id := results[0].GetUnsafeStringOrEmpty()
		bidId := results[1].GetUnsafeStringOrEmpty()
		bidImpId := results[2].GetUnsafeStringOrEmpty()
		bidPrice := results[3].GetUnsafeStringOrEmpty()
		bidAdm := results[4].GetUnsafeStringOrEmpty()
		cur := results[5].GetUnsafeStringOrEmpty()
		foo := results[6].GetUnsafeStringOrEmpty()

		if "9f3701ef-1d8a-4b67-a601-e9a1a2fbcb7c" != id ||
			"0" != bidId ||
			"1" != bidImpId ||
			"0.35258" != bidPrice ||
			`some \"html\" code\nnext line` != bidAdm ||
			"USD" != cur ||
			"" != foo {
			panic("ack")
		}
	}
}
*/

/*
BenchmarkEachKey
BenchmarkEachKey-8
  308006	      3885 ns/op	     272 B/op	       4 allocs/op
PASS
*/
func BenchmarkEachKey(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id string
		var bidId string
		var bidImpId string
		var bidPrice string
		var bidAdm string
		var cur string
		var foo string

		bodyBytes := StringAsBytes(BidResponse_Banner)

		jsonparser.EachKey(bodyBytes, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
			switch idx {
			case 0:
				id = BytesAsString(value)
			case 1:
				bidId = BytesAsString(value)
			case 2:
				bidImpId = BytesAsString(value)
			case 3:
				bidPrice = BytesAsString(value)
			case 4:
				bidAdm = BytesAsString(value)
			case 5:
				cur = BytesAsString(value)
			case 6:
				foo = BytesAsString(value)
			}
		},
			[][]string{
				{"id"},
				{"seatbid", "[0]", "bid", "[0]", "id"},
				{"seatbid", "[0]", "bid", "[0]", "impid"},
				{"seatbid", "[0]", "bid", "[0]", "price"},
				{"seatbid", "[0]", "bid", "[0]", "adm"},
				{"cur"},
				{"foo"},
			}...,
		)

		if "9f3701ef-1d8a-4b67-a601-e9a1a2fbcb7c" != id ||
			"0" != bidId ||
			"1" != bidImpId ||
			"0.35258" != bidPrice ||
			`some \"html\" code\nnext line` != bidAdm ||
			"USD" != cur ||
			"" != foo {
			panic("ack")
		}
	}

}

const BidResponse_Banner = `
{
	"id":"9f3701ef-1d8a-4b67-a601-e9a1a2fbcb7c",
	"seatbid":[
	   {
		  "bid":[
			 {
				"id":"0",
				"adomain":[
				   "someurl.com",
				   "someotherurl.com"
				],
				"impid":"1",
				"price":0.35258,
				"adm":"some \"html\" code\nnext line",
				"w":300,
				"h":250
			 }
		  ]
	   }
	],
	"cur":"USD"
}
`

var BidResponse_Banner_Paths [][]string = [][]string{
	{"id"},
	{"seatbid", "[0]", "bid", "[0]", "id"},
	{"seatbid", "[0]", "bid", "[0]", "impid"},
	{"seatbid", "[0]", "bid", "[0]", "price"},
	{"seatbid", "[0]", "bid", "[0]", "adm"},
	{"cur"},
	{"foo"},
}
