package goj

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	input := `{"hi":"mom"}`
	output :=
		`{
  "hi": "mom"
}`

	r := strings.NewReader(input)
	dec := NewDecoder(r)

	if err := dec.Decode(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, output, dec.String())
}
