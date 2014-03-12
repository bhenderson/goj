package goj

import (
	"fmt"
	"strconv"
	"strings"
)

type stateFunc func(*Path, string) (stateFunc, int)

type SyntaxError struct {
	msg string
}

func (e *SyntaxError) Error() string {
	return e.msg
}

func (p *Path) parse() error {
	f := stateKey
	data := p.p
	var i, x int

	for f != nil && i <= len(data) {
		f, x = f(p, data[i:])
		i = i + x
	}

	return p.err
}

// beginning state
func stateKey(p *Path, str string) (f stateFunc, i int) {
	data := []byte(str)
L:
	for ; i < len(data); i++ {
		switch data[i] {
		case '\\':
			// escape char. skip next value
			i++
			if i >= len(data) {
				f = stateEscape
				return
			}
			data = append(data[:i-1], data[i:]...)
		case '=':
			f = stateValue
			if i == 0 {
				p.sel = append(p.sel, "*")
			}
			break L
		case '[':
			f = stateArray
			break L
		case '.':
			f = stateSep
			break L
		}
	}
	if i != 0 {
		p.sel = append(p.sel, string(data[:i]))
	}
	i++
	return
}

// state after key=
func stateValue(p *Path, data string) (f stateFunc, i int) {
	var pair Pair
	x := len(p.sel) - 1
	pair.key = p.sel[x].(string)
	p.sel[x] = pair

	return stateValueNum, i
}

func setValue(p *Path, v interface{}) {
	i := len(p.sel) - 1
	pair := p.sel[i].(Pair)
	pair.val = v
	p.sel[i] = pair
}

// state after key= (try number)
func stateValueNum(p *Path, data string) (f stateFunc, i int) {
	var fl float32
	var st string

	r := strings.NewReader(data)
	x, _ := fmt.Fscanf(r, "%f%s", &fl, &st)

	if x > 0 {
		setValue(p, fl)
		i = len(data) - len(st)
		if data[i] == '.' {
			f = stateParent
		}
	} else {
		f = stateValueString
	}

	return
}

// state after key=\D
func stateValueString(p *Path, data string) (f stateFunc, i int) {
L:
	for ; i < len(data); i++ {
		switch data[i] {
		case '\\':
			i++
			if i >= len(data) {
				f = stateEscape
				return
			}
		case '.':
			f = stateParent
			break L
		}
	}

	if i != 0 {
		setValue(p, data[:i])
	}

	i++

	return
}

// state after key[
func stateArray(p *Path, data string) (f stateFunc, i int) {
L:
	for ; i < len(data); i++ {
		switch data[i] {
		case ']':
			break L
		}
	}

	x, e := strconv.Atoi(data[:i])
	if e != nil {
		f = addError(`invalid index`)
		return
	}
	p.sel = append(p.sel, x)

	f = stateArrayEnd
	i++

	return
}

// state after ]
func stateArrayEnd(p *Path, data string) (f stateFunc, i int) {
	if len(data) == 0 {
		return
	}

	if data[i] == '.' {
		f = stateSep
	} else {
		f = addError(`expected "."`)
	}
	i++

	return
}

// state after .
func stateSep(p *Path, data string) (f stateFunc, i int) {
	if len(data) > 0 && data[0] == '.' {
		return stateParent, i
	}
	return stateKey, i
}

// state after **
func stateRecursive(p *Path, data string) (f stateFunc, i int) {
	return
}

// state after ..
func stateParent(p *Path, data string) (f stateFunc, i int) {
	if !(len(data) > 0 && data[0] == '.') {
		f = addError(`expected ".."`)
		return
	}
	p.sel = append(p.sel, "..")
	i++
	return stateKey, i
}

// addError returns a stateFunc which sets the error.
func addError(msg string) stateFunc {
	return func(p *Path, data string) (f stateFunc, i int) {
		p.err = &SyntaxError{"invalid path at " + p.p[:len(p.p)-len(data)] + " " + msg}
		return
	}
}

var stateEscape = addError("invalid escape character")
