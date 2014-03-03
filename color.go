package goj

import (
	"bytes"
	"encoding/json"
	"github.com/mgutz/ansi"
)

var reset = []byte(ansi.ColorCode("reset"))

type colorFunc func(dst *bytes.Buffer, v interface{}) error

func makeColorFunc(style string) colorFunc {
	color := []byte(ansi.ColorCode(style))
	return func(dst *bytes.Buffer, v interface{}) (err error) {
		b, err := json.Marshal(v)
		if err != nil {
			return
		}
		dst.Write(color)
		dst.Write(b)
		dst.Write(reset)
		return
	}
}

var (
	blue   = makeColorFunc("blue")
	green  = makeColorFunc("green")
	grey   = makeColorFunc("black+b")
	yellow = makeColorFunc("yellow")
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
			err = blue(buf, k)
			if err != nil {
				return err
			}
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
		err = yellow(buf, x)
		if err != nil {
			return err
		}
	case string:
		green(buf, x)
		if err != nil {
			return err
		}
	case nil:
		grey(buf, x)
		if err != nil {
			return err
		}
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
