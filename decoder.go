package goj

import (
	"bytes"
	"encoding/json"
	"os"
)

// io.Reader without io
type reader interface {
	Read(p []byte) (n int, err error)
}

func NewDecoder(f *os.File) (d *Decoder) {
	d = &Decoder{file: f, dec: json.NewDecoder(f)}
	return
}

// BUG(bh) need public method to access Decoder.v
type Decoder struct {
	color colorSet
	dec   *json.Decoder
	file  *os.File
	v     interface{}
	ind   *indent
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

func (d *Decoder) FileName() string {
	return d.file.Name()
}

// SetColor sets the option to colorize the pretty formatting. Takes one of Colors.
func (d *Decoder) SetColor(set colorSet) {
	d.color = set
}

// String returns nicely formatted json, optionally colored.
func (d *Decoder) String() string {
	id := d.indent()

	if d.color.IsTrue() {
		var buf bytes.Buffer
		colorize(&buf, d.v, id)

		return buf.String()
	}

	return d.StringColorless()
}

func (d *Decoder) StringColorless() string {
	// TODO move this into color
	id := d.indent()

	b, err := json.MarshalIndent(d.v, id.prefix, id.indent)

	// TODO better error handling.
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (d *Decoder) indent() *indent {
	if d.ind == nil {
		d.ind = &indent{indent: "  "}
	}

	return d.ind
}
