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
	plain  = makeColorFunc("plain")
)

type indent struct {
	prefix, indent string
	depth          int
}

func (idt *indent) start(buf *bytes.Buffer, b byte) {
	idt.depth++
	buf.WriteByte(b)
}

func (idt *indent) end(buf *bytes.Buffer, b byte, p int) {
	idt.depth--
	if p != 0 {
		buf.Truncate(p)
		newline(buf, idt)
	}
	buf.WriteByte(b)
}

func colorize(buf *bytes.Buffer, v interface{}, idt *indent) (err error) {
	var p int

	switch x := v.(type) {
	case map[string]interface{}:
		// json keys must be strings
		idt.start(buf, '{')
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
		idt.end(buf, '}', p)
	case []interface{}:
		idt.start(buf, '[')
		for _, val := range x {
			newline(buf, idt)
			err = colorize(buf, val, idt)
			if err != nil {
				return err
			}
			p = buf.Len()
			buf.WriteByte(',')
		}
		idt.end(buf, ']', p)
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
		err = plain(buf, x)
		if err != nil {
			return err
		}
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
