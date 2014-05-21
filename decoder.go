package goj

import (
	"encoding/json"
)

// io.Reader without io
type reader interface {
	Read(p []byte) (n int, err error)
}

// os.File
type File interface {
	reader
	Name() string
}

func NewDecoder(files ...File) (d *Decoder) {
	//d = &Decoder{files: files, dec: json.NewDecoder(f)}
	out := make(chan *Val)
	d = &Decoder{files: files, outc: out}
	go internDecode(d)
	return
}

type Decoder struct {
	files []File
	outc  chan *Val
	ind   *indent
	color colorSet
}

// Decode takes a filter string and decodes from reader.
func (d *Decoder) Decode(f string) <-chan *Val {
	out := make(chan *Val)
	go func() {
		for v := range d.outc {
			filterOn(v, f)
			out <- v
		}
		close(out)
	}()
	return out
}

// internal decode loops through all files
func internDecode(d *Decoder) {
	for _, f := range d.files {
		dec := json.NewDecoder(f)
		for {
			var v interface{}
			if err := dec.Decode(&v); err != nil {
				if err.Error() == "EOF" {
					break
				} else {
					// TODO handle error
					d.outc <- &Val{Error: err}
					break
				}
			} else {
				val := &Val{v: v, file: f, dec: d}
				d.outc <- val
			}
		}
	}
	close(d.outc)
}

// SetColor sets the option to colorize the pretty formatting. Takes one of Colors.
func (d *Decoder) SetColor(set colorSet) {
	d.color = set
}

// should be on the option object
func (d *Decoder) indent() *indent {
	if d.ind == nil {
		d.ind = defaultIndent
	}

	return d.ind
}
