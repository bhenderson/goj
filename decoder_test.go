package goj

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_DecodeMultipleInputs(t *testing.T) {
	input := `
		{"a":"b"}
		{"c":"d"}
	`
	dec := testDecoder(t, input)

	act := <-dec.outc
	exp := `{"a":"b"}`
	assert.Equal(t, testMarshal(t, exp), act.v)

	act = <-dec.outc
	exp = `{"c":"d"}`
	assert.Equal(t, testMarshal(t, exp), act.v)

	act = <-dec.outc
	assert.Equal(t, (*Val)(nil), act)
}

func ExampleNewDecoder() {
	// Decode a line of json at a time, optionally filtering the result.
	filter := ""
	files := []File{os.Stdin}

	dec := NewDecoder(files...)
	dec.SetColor(ColorAlways)

	for val := range dec.Decode(filter, false) {
		fmt.Println(val)
	}
}
