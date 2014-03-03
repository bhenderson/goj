package goj

import (
	"bytes"
	"encoding/json"
	"github.com/mgutz/ansi"
)

var (
	blue   = ansi.ColorFunc("blue")
	green  = ansi.ColorFunc("green")
	grey   = ansi.ColorFunc("black+b")
	yellow = ansi.ColorFunc("yellow")
)

func (d Decoder) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := colorize(&buf, d.v)
	return buf.Bytes(), err
}

func colorize(buf *bytes.Buffer, v interface{}) (err error) {
	var b []byte

	switch x := v.(type) {
	case map[string]interface{}:
		// json keys must be strings
		buf.WriteRune('{')
		for k, val := range x {
			b, err = json.Marshal(k)
			if err != nil {
				return err
			}
			buf.WriteString(blue(string(b)))
			buf.WriteRune(':')
			err = colorize(buf, val)
			if err != nil {
				return err
			}
			buf.WriteRune(',')
		}
		buf.Truncate(buf.Len() - 1) // last ,
		buf.WriteRune('}')
	case []interface{}:
		buf.WriteRune('[')
		for _, val := range x {
			err = colorize(buf, val)
			if err != nil {
				return err
			}
			buf.WriteRune(',')
		}
		buf.Truncate(buf.Len() - 1) // last ,
		buf.WriteRune(']')
	case int, float64:
		b, err = json.Marshal(x)
		if err != nil {
			return err
		}
		buf.WriteString(yellow(string(b)))
	case string:
		b, err = json.Marshal(x)
		if err != nil {
			return err
		}
		buf.WriteString(green(string(b)))
	case nil:
		b, err = json.Marshal(x)
		if err != nil {
			return err
		}
		buf.WriteString(grey(string(b)))
	default:
		b, err = json.Marshal(x)
		if err != nil {
			return err
		}
		buf.Write(b)
	}
	return
}
