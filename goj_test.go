package goj

import (
	"encoding/json"
	"strings"
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

type intTest interface {
	Fatal(args ...interface{})
}

func testVal(t intTest, input string) *Val {
	d := testDecoder(t, input)
	return <-d.Decode("")
}

func testDecoder(t intTest, input string) *Decoder {
	f := &testFile{
		strings.NewReader(input),
		"test input",
	}
	return NewDecoder(f)
}

func testMarshal(t intTest, input string) interface{} {
	r := strings.NewReader(input)
	dec := json.NewDecoder(r)

	var v interface{}
	if err := dec.Decode(&v); err != nil {
		t.Fatal(err)
	}
	return v
}
