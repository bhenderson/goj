package goj

import (
	"bytes"
	"encoding/json"
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
	dec *json.Decoder
	v   interface{}
}

func (d *Decoder) Decode(f string) (err error) {
	err = d.dec.Decode(&d.v)
	if err != nil {
		return
	}
	if f != "" {
		d.FilterOn(f)
	}
	return
}

func (d *Decoder) String() string {
	var buf bytes.Buffer
	colorize(&buf, d.v, &indent{indent: "  "})

	return buf.String()
}
