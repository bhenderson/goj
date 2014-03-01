package goj

import (
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

func (d *Decoder) String() string {
	b, _ := json.MarshalIndent(d.v, "", "  ")

	return string(b)
}
