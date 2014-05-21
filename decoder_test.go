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
	out := dec.Decode("")

	act := <-out
	exp := `{"a":"b"}`
	assert.Equal(t, testMarshal(t, exp), act.v)

	act = <-out
	exp = `{"c":"d"}`
	assert.Equal(t, testMarshal(t, exp), act.v)

	act = <-out
	assert.True(t, nil == act)
}

func TestDecodeInvalid(t *testing.T) {
	v := <-testDecoder(t, `{"invalid`).Decode("")

	assert.Nil(t, (&Val{}).Error)
	assert.NotNil(t, v.Error)
}

func TestDecodeNumber(t *testing.T) {
	input := `{"foo":1111111111222222222233333333334444444444}`

	dec := testDecoder(t, input)
	act := <-dec.Decode("")

	data := "{\n  \"foo\": 1111111111222222222233333333334444444444\n}"

	assert.Equal(t, data, act.String())
}

func ExampleNewDecoder() {
	// Decode a line of json at a time, optionally filtering the result.
	filter := ""
	files := []File{os.Stdin}

	dec := NewDecoder(files...)
	dec.SetColor(ColorAlways)

	for val := range dec.Decode(filter) {
		fmt.Println(val)
	}
}
