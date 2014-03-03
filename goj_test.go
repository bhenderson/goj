package goj

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_Decode(t *testing.T) {
	input := `{
		"hi":"mom",
		"n":null,
		"a":["b",1]
	}`
	output :=
		`{
  "\x1b[34ma\x1b[0m": [
  	"\x1b[32mb\x1b[0m",
  	"\x1b[33m1\x1b[0m"
  ],
  "\x1b[34mhi\x1b[0m": "\x1b[32mmom\x1b[0m",
  "\x1b[34mn\x1b[0m": \x1b[1;30mnull\x1b[0m
}`

	r := strings.NewReader(input)
	dec := NewDecoder(r)

	if err := dec.Decode(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, output, dec.String())
}
