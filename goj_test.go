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

	output := `{
  [34m"hi"[0m: [32m"mom"[0m,
  [34m"n"[0m: [1;30mnull[0m,
  [34m"a"[0m: [
    [32m"b"[0m,
    [33m1[0m
  ]
}`

	r := strings.NewReader(input)
	dec := NewDecoder(r)

	if err := dec.Decode(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, output, dec.String())
}
