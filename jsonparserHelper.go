package jsonparserHelper

import (
	"bytes"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/buger/jsonparser"
	"github.com/juju/errors"
)

// manually hide an allocation from escape analysis, but go vet/lint doesnt like this
/*
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
*/

func BytesAsString(bs []byte) string {
	/* same as:
	from strings.Builder
	func (b *Builder) String() string {
		return *(*string)(unsafe.Pointer(&b.buf))
	}
	*/
	return *(*string)(unsafe.Pointer(&bs))
}

func StringAsBytes(s string) []byte {
	// I dont think this is safe
	//return *(*[]byte)(unsafe.Pointer(&s))

	// this should be safe
	strh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	var sh reflect.SliceHeader
	sh.Data = strh.Data
	sh.Len = strh.Len
	sh.Cap = strh.Len
	return *(*[]byte)(unsafe.Pointer(&sh))
}

type JsonparserResult struct {
	Key   string
	Idx   int
	Value []byte
	Vt    jsonparser.ValueType
	Err   error
}

func MakeJsonparserResult(paths [][]string, idx int, value []byte, vt jsonparser.ValueType, err error) JsonparserResult {
	path := paths[idx]
	return JsonparserResult{
		Key:   path[len(path)-1],
		Idx:   idx,
		Value: value,
		Vt:    vt,
		Err:   err,
	}
}

func (t JsonparserResult) GetUnsafeStringOrEmpty() string {
	if t.Vt == jsonparser.String || t.Vt == jsonparser.Number || t.Vt == jsonparser.Boolean {
		return BytesAsString(t.Value)
	}
	return ""
}

var trueValue []byte = []byte("true")

func (t JsonparserResult) GetInt64OrZero() int64 {
	if bytes.Equal(t.Value, trueValue) {
		return 1
	}

	if t.Vt == jsonparser.String || t.Vt == jsonparser.Number {
		i, err := strconv.ParseInt(BytesAsString(t.Value), 10, 64)
		if err == nil {
			return i
		}
		return 0
	}
	return 0
}

func (t JsonparserResult) GetFloatOrZero() float64 {
	if bytes.Equal(t.Value, trueValue) {
		return 1
	}

	if t.Vt == jsonparser.String || t.Vt == jsonparser.Number {
		i, err := strconv.ParseFloat(BytesAsString(t.Value), 64)
		if err == nil {
			return i
		}
		return 0
	}
	return 0
}

func (t JsonparserResult) GetIntOrZero() int {
	return int(t.GetInt64OrZero())
}

const JsonparserValueType_UnsafeString = 1
const JsonparserValueType_Int64 = 2
const JsonparserValueType_Int = 3
const JsonparserValueType_Float64 = 4

type JsonparserValue struct {
	Type   int
	Offset uintptr
}

func UnsafeStringValue(p uintptr) JsonparserValue {
	return JsonparserValue{Type: JsonparserValueType_UnsafeString, Offset: p}
}

func Int64Value(p uintptr) JsonparserValue {
	return JsonparserValue{Type: JsonparserValueType_Int64, Offset: p}
}

func IntValue(p uintptr) JsonparserValue {
	return JsonparserValue{Type: JsonparserValueType_Int, Offset: p}
}

func Float64Value(p uintptr) JsonparserValue {
	return JsonparserValue{Type: JsonparserValueType_Float64, Offset: p}
}

type Parser struct {
	Path  []string
	Value JsonparserValue
}

type Parsers []Parser

// parser = append(parser, Parser{[]string{"foo"}, UnsafeStringValue(unsafe.Offsetof(Fubar{}.Foo))})
func MakeParser(storage []Parser) Parsers {
	return storage
}

func (t Parsers) Add(v JsonparserValue, keys ...string) Parsers {
	return append(t, Parser{keys, v})
}

func RunParser(bodyBytes []byte, parser []Parser, dest unsafe.Pointer) error {
	var e error
	paths := make([][]string, 128)[:0]
	for i := 0; i < len(parser); i++ {
		paths = append(paths, parser[i].Path)
	}

	jsonparser.EachKey(bodyBytes, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
		if err != nil {
			if e == nil {
				e = err
			} else {
				e = errors.Wrap(e, err)
			}
			return
		}

		r := MakeJsonparserResult(paths, idx, value, vt, err)

		parseValue := parser[idx].Value
		pValue := unsafe.Pointer(uintptr(dest) + parseValue.Offset)
		switch parseValue.Type {
		case JsonparserValueType_UnsafeString:
			*(*string)(pValue) = r.GetUnsafeStringOrEmpty()
		case JsonparserValueType_Int64:
			*(*int64)(pValue) = r.GetInt64OrZero()
		case JsonparserValueType_Int:
			*(*int)(pValue) = r.GetIntOrZero()
		case JsonparserValueType_Float64:
			*(*float64)(pValue) = r.GetFloatOrZero()
		}

	}, paths...)

	return e
}

/*

type Fubar struct {
	Foo int
	Bar int
}

fubar := &fubar{ foo:1, bar:2 }
err := Parse(bytes, fubar, [
	{ "foo", StringValue(&fubar.Foo) },
	{ "bar", Int64Value(&fubar.Bar) },
])

*/
