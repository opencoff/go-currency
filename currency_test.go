package currency_test

import (
	"encoding/json"
	"runtime"
	"testing"

	"currency"
)

func assert(cond bool, t *testing.T) {

	if cond {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	t.Fatalf("%s: %d: Assertion failed\n", file, line)
}

type testcase struct {
	in    string
	out   string
	oprec int
}

var tests = [...]testcase{
	{"123456.987654321", "123456.987654321000", 12},
	{"123456.987654321", "123456.98765", 5},
	{"123456.987654321", "123456.987654", 6},
	{"123.00456000000", "123.00456", 5},
	{"123.00456000000", "123.004560", 6},
	{"000.00456000000", "0.004", 3},
	{"000.00456000000", "0.00456", 5},
	{"123.00000", "123.0", 1},
	{"123.00000", "123.00", 2},
	{"000.0000", "0.0", 1},
}

func Test_fmt(t *testing.T) {
	for _, tc := range tests {
		c, err := currency.NewFromString(tc.in)
		assert(err == nil, t)

		s := c.StringFixed(tc.oprec)
		t.Logf("in=|%s|, exp out=|%s| => |%s|\n", tc.in, tc.out, s)
		assert(s == tc.out, t)
	}
}

func Test_json(t *testing.T) {
	ii, err := currency.NewFromString("123.0005430123")
	assert(err == nil, t)

	m, err := json.Marshal(ii)
	assert(err == nil, t)

	t.Logf("json=%s\n", string(m))

	var xx currency.Currency

	err = json.Unmarshal(m, &xx)
	assert(err == nil, t)
	assert(ii.Eq(&xx), t)
}

func Benchmark_NewFromString(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_, _ = currency.NewFromString("0023.0045600000")
	}
}

func Benchmark_String0(b *testing.B) {
	ii, err := currency.NewFromString("0023.0045600000")
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_ = ii.String()
	}
}

func Benchmark_String1(b *testing.B) {
	ii, err := currency.NewFromString("0023.0000004567")
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		_ = ii.String()
	}
}
