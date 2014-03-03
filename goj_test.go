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
		"{\n  \x1b[34m\"hi\"\x1b[0m: \x1b[32m\"mom\"\x1b[0m,\n  \x1b[34m\"n\"\x1b[0m: \x1b[1;30mnull\x1b[0m,\n  \x1b[34m\"a\"\x1b[0m: [\n    \x1b[32m\"b\"\x1b[0m,\n    \x1b[33m1\x1b[0m\n  ]\n}"

	r := strings.NewReader(input)
	dec := NewDecoder(r)

	if err := dec.Decode(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, output, dec.String())
}
