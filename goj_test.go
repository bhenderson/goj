package goj

import (
	"encoding/json"
	"os"
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

type intTest interface {
	Fatal(args ...interface{})
}

func testDecoder(t intTest, input string) *Decoder {
	r := strings.NewReader(input)
	f := os.Stdin
	dec := &Decoder{file: f, dec: json.NewDecoder(r)}

	if err := dec.Decode(""); err != nil {
		t.Fatal(err)
	}
	return dec
}
