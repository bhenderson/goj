package goj

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func Test_DecodeMultipleInputs(t *testing.T) {
	input := `
		{"a":"b"}
		{"c":"d"}
	`
	dec := testDecoder(t, input)

	exp := `{"a":"b"}`
	assert.Equal(t, testDecoder(t, exp).v, dec.v)

	if err := dec.Decode(); err != nil {
		t.Fatal(err)
	}
	exp = `{"c":"d"}`
	assert.Equal(t, testDecoder(t, exp).v, dec.v)

	if err := dec.Decode(); err != io.EOF {
		t.Fatal("expected EOF, got", err)
	}
}

// helpers

func testDecoder(t *testing.T, input string) *Decoder {
	r := strings.NewReader(input)
	dec := NewDecoder(r)

	if err := dec.Decode(); err != nil {
		t.Fatal(err)
	}
	return dec
}
