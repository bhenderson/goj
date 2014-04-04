package goj

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
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
	// Read a line of stding at a time, parsing it as json, then filtering the result.
	// dec.String() will return pretty printed json as a string.
	filter := ""
	reader := strings.NewReader(`{"hi":"mom"}{"foo":"bar"}`)
	dec := NewDecoder(reader)

	for {
		if err := dec.Decode(filter); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println(dec)
	}
	// Output: {
	//   "hi": "mom"
	//}
	//{
	//   "foo": "bar"
	//}
}
