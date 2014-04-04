package goj

import (
	"bytes"
	"encoding/json"
	"github.com/bhenderson/terminal"
)

// io.Reader without io
type reader interface {
	Read(p []byte) (n int, err error)
}

func NewDecoder(r reader) (d *Decoder) {
	d = &Decoder{dec: json.NewDecoder(r)}
	return
}

// BUG(bh) need public method to access Decoder.v
type Decoder struct {
	dec   *json.Decoder
	v     interface{}
	color colorSet
}

// Val is the attribute reader for getting the decoded json value.
func (d *Decoder) Val() interface{} {
	return d.v
}

// Decode takes a filter string and decodes from reader.
func (d *Decoder) Decode(f string) (err error) {
	err = d.dec.Decode(&d.v)
	if err == nil && f != "" {
		err = filterOn(d, f)
	}
	return
}

// SetColor sets the option to colorize the pretty formatting. Takes one of Colors.
func (d *Decoder) SetColor(set colorSet) {
	d.color = set
}

// String returns nicely formatted json, optionally colored.
func (d *Decoder) String() string {
	id := &indent{indent: "  "}

	if shouldColor(d.color) {
		var buf bytes.Buffer
		colorize(&buf, d.v, id)

		return buf.String()
	}

	// TODO move this into color
	b, err := json.MarshalIndent(d.v, id.prefix, id.indent)

	// TODO better error handling.
	if err != nil {
		panic(err)
	}

	return string(b)
}

func shouldColor(set colorSet) (b bool) {
	switch set {
	case ColorAlways:
		b = true
	case ColorNever:
		b = false
	case ColorAuto:
		b = terminal.IsTerminal(1)
	}

	return b
}
