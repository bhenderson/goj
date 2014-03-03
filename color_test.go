package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_colorize(t *testing.T) {
	input := map[string]interface{}{
		"hi":  "mom",
		"foo": nil,
		"a": []interface{}{
			"b",
			1,
		},
	}
	expected :=
		"{\x1b[34m\"hi\"\x1b[0m:\x1b[32m\"mom\"\x1b[0m,\x1b[34m\"foo\"\x1b[0m:\x1b[1;30mnull\x1b[0m,\x1b[34m\"a\"\x1b[0m:[\x1b[32m\"b\"\x1b[0m,\x1b[33m1\x1b[0m]}"

	b, err := Decoder{v: input}.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, string(b))
}
