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
		if len(color) != 0 {
			dst.Write(color)
		}
		dst.Write(b)
		if len(color) != 0 {
			dst.Write(reset)
		}
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

func colorize(buf *bytes.Buffer, v interface{}, idt *indent) (err error) {
	var p int

	switch x := v.(type) {
	case map[string]interface{}:
		// json keys must be strings
		prefix(idt, buf, '{')
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
		postfix(idt, buf, '}', p)
	case []interface{}:
		prefix(idt, buf, '[')
		for _, val := range x {
			newline(buf, idt)
			err = colorize(buf, val, idt)
			if err != nil {
				return err
			}
			p = buf.Len()
			buf.WriteByte(',')
		}
		postfix(idt, buf, ']', p)
	case int, float64:
		err = yellow(buf, x)
	case string:
		err = green(buf, x)
	case nil:
		err = grey(buf, x)
	default:
		err = plain(buf, x)
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

func prefix(idt *indent, buf *bytes.Buffer, b byte) {
	idt.depth++
	buf.WriteByte(b)
}

func postfix(idt *indent, buf *bytes.Buffer, b byte, p int) {
	idt.depth--
	if p != 0 {
		buf.Truncate(p)
		newline(buf, idt)
	}
	buf.WriteByte(b)
}
