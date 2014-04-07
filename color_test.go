package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Decode(t *testing.T) {
	input := `{
		"hi":"mom",
		"n":null,
		"a":["b",1]
	}`

	output := `{
  [34m"hi"[0m: [32m"mom"[0m,
  [34m"n"[0m: [1;30mnull[0m,
  [34m"a"[0m: [
    [32m"b"[0m,
    [33m1[0m
  ]
}`

	dec := testDecoder(t, input)
	dec.SetColor(ColorAuto)

	assert.Equal(t, output, dec.String())
}

func TestDecode_noColor(t *testing.T) {
	exp := `{
  "a": [
    "b",
    1
  ],
  "hi": "mom",
  "n": null
}`

	d := testDecoder(t, exp)

	d.Decode("")

	assert.Equal(t, exp, d.String())
}
