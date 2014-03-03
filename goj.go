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

func (d *Decoder) Decode() (err error) {
	return d.dec.Decode(&d.v)
}

func (d Decoder) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := colorize(&buf, d.v, &indent{indent: "  "})
	return buf.Bytes(), err
}

func (d *Decoder) String() string {
	b, _ := d.MarshalJSON()

	return string(b)
}
