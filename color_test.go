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

	v := testVal(t, input)
	v.dec.SetColor(ColorAuto)

	assert.Equal(t, output, v.String())
}

func TestDecode_noColor(t *testing.T) {
	output := `{
  "a": [
    "b",
    1
  ],
  "hi": "mom",
  "n": null
}`

	v := testVal(t, output)
	v.dec.SetColor(ColorNever)

	assert.Equal(t, output, v.String())
}
