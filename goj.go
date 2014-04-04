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

type Decoder struct {
	dec   *json.Decoder
	v     interface{}
	color *ColorSet
}

func (d *Decoder) Decode(f string) (err error) {
	err = d.dec.Decode(&d.v)
	if err == nil && f != "" {
		err = d.FilterOn(f)
	}
	return
}

func (d *Decoder) SetColor(set *ColorSet) {
	d.color = set
}

func (d *Decoder) String() string {
	id := &indent{indent: "  "}

	if shouldColor(d.color) {
		var buf bytes.Buffer
		colorize(&buf, d.v, id)

		return buf.String()
	}

	b, err := json.MarshalIndent(d.v, id.prefix, id.indent)

	if err != nil {
		panic(err)
	}

	return string(b)
}

func shouldColor(set *ColorSet) (b bool) {
	switch set.c {
	case colorAlways:
		b = true
	case colorNever:
		b = false
	case colorAuto:
		b = terminal.IsTerminal(1)
	}

	return b
}
