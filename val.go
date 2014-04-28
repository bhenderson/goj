package goj

import (
	"bytes"
	"encoding/json"
)

type Val struct {
	v    interface{}
	file File
	dec  *Decoder // should be options, not Decoder
}

func (v *Val) FileName() string {
	return v.file.Name()
}

// String returns nicely formatted json or a diff, optionally colored.
func (v *Val) String() string {
	id := v.dec.indent()

	if v.dec.color.IsTrue() {
		var buf bytes.Buffer
		colorize(&buf, v.v, id)

		return buf.String()
	}

	b, _ := v.MarshalJSON()
	return string(b)
}

func (v *Val) MarshalJSON() ([]byte, error) {
	// TODO move this into color
	id := v.dec.indent()

	return json.MarshalIndent(v.v, id.prefix, id.indent)
}
