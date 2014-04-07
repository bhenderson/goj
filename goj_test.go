package goj

import (
	"encoding/json"
	"os"
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

func testDecoder(t *testing.T, input string) *Decoder {
	r := strings.NewReader(input)
	f := os.Stdin
	dec := &Decoder{file: f, dec: json.NewDecoder(r)}

	if err := dec.Decode(""); err != nil {
		t.Fatal(err)
	}
	return dec
}
