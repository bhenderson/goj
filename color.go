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

type indent struct {
	prefix, indent string
	depth          int
}

func (d Decoder) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := colorize(&buf, d.v, &indent{indent: "  "})
	return buf.Bytes(), err
}

func colorize(buf *bytes.Buffer, v interface{}, idt *indent) (err error) {
	var b []byte
	var p int

	switch x := v.(type) {
	case map[string]interface{}:
		// json keys must be strings
		idt.depth++
		buf.WriteByte('{')
		for k, val := range x {
			newline(buf, idt)
			b, err = json.Marshal(k)
			if err != nil {
				return err
			}
			buf.WriteString(blue(string(b)))
			buf.WriteByte(':')
			buf.WriteByte(' ')
			err = colorize(buf, val, idt)
			if err != nil {
				return err
			}
			p = buf.Len()
			buf.WriteByte(',')
		}
		idt.depth--
		if p != 0 {
			buf.Truncate(buf.Len() - 1) // last ,
			newline(buf, idt)
		}
		buf.WriteByte('}')
	case []interface{}:
		idt.depth++
		buf.WriteByte('[')
		for _, val := range x {
			newline(buf, idt)
			err = colorize(buf, val, idt)
			if err != nil {
				return err
			}
			p = buf.Len()
			buf.WriteByte(',')
		}
		idt.depth--
		if p != 0 {
			buf.Truncate(buf.Len() - 1) // last ,
			newline(buf, idt)
		}
		buf.WriteByte(']')
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

func newline(buf *bytes.Buffer, idt *indent) {
	buf.WriteByte('\n')
	buf.WriteString(idt.prefix)
	for i := 0; i < idt.depth; i++ {
		buf.WriteString(idt.indent)
	}
}
