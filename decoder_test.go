package goj

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
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

	if err := dec.Decode(""); err != nil {
		t.Fatal(err)
	}
	exp = `{"c":"d"}`
	assert.Equal(t, testDecoder(t, exp).v, dec.v)

	if err := dec.Decode(""); err != io.EOF {
		t.Fatal("expected EOF, got", err)
	}
}

func ExampleNewDecoder() {
	// Decode a line of json at a time, optionally filtering the result.
	filter := ""
	dec := NewDecoder(os.Stdin)

	for {
		if err := dec.Decode(filter); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// dec.Val() -> decoded value
	}
}
