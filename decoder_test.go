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
	out := dec.Decode("", false)

	act := <-out
	exp := `{"a":"b"}`
	assert.Equal(t, testMarshal(t, exp), act.v)

	act = <-out
	exp = `{"c":"d"}`
	assert.Equal(t, testMarshal(t, exp), act.v)

	act = <-out
	assert.True(t, nil == act)
}

func Test_DecodeDiffSingleInput(t *testing.T) {
	input := `
		{"a":"b"}
	`
	dec := testDecoder(t, input)
	out := dec.Decode("", true)

	act := <-out
	exp := `{"a":"b"}`
	assert.Equal(t, testMarshal(t, exp), act.v)
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
