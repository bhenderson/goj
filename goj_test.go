package goj

import (
	"encoding/json"
	"strings"
	"testing"
)

// test helpers

var input = `{
	"store": {
		"bicycles": [
			{
				"color": "red",
				"price": 3.99
			},
			{
				"color": "blue",
				"price": 2.99
			}
		],
		"truck": {
			"color": "yellow",
			"price": 3.99
		}
	}
}`

type testFile struct {
	r *strings.Reader
	n string
}

func (t *testFile) Read(p []byte) (int, error) {
	return t.r.Read(p)
}

func (t *testFile) Name() string {
	return t.n
}

func testVal(t testing.TB, input string) *Val {
	d := testDecoder(t, input)
	return <-d.Decode("")
}

func testDecoder(t testing.TB, input string) *Decoder {
	f := &testFile{
		strings.NewReader(input),
		"test input",
	}
	return NewDecoder(f)
}

func testMarshal(t testing.TB, input string) interface{} {
	r := strings.NewReader(input)
	dec := json.NewDecoder(r)
	// TODO this logic should be shared with the lib.
	dec.UseNumber()

	var v interface{}
	if err := dec.Decode(&v); err != nil {
		t.Fatal(err)
	}
	return v
}
