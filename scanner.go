package goj

import (
	"fmt"
	"strings"
)

type stateFunc func(*Path, string) (stateFunc, int)

type syntaxError struct {
	msg string
}

func (e *syntaxError) Error() string {
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

func appendPathSel(p *Path, s pathSel) {
	p.sel = append(p.sel, s)
}

// beginning state
func stateKey(p *Path, data string) (f stateFunc, i int) {
	var escapes int
L:
	for ; i < len(data); i++ {
		switch data[i] {
		case '\\':
			i++
			// escape char. skip next value
			if i >= len(data) {
				f = stateEscape
				return
			}
			data = data[:i-1] + data[i:]
			i--
			escapes++
		case '=':
			f = stateValue
			if i == 0 {
				appendPathSel(p, &pathKey{"*"})
			}
			break L
		case '[':
			f = stateArray
			break L
		case '.':
			f = stateSep
			break L
		case '*':
			if len(data) > i+1 && data[i+1] == '*' {
				return stateRecursive, i
			}
		}
	}
	if i != 0 {
		appendPathSel(p, &pathKey{string(data[:i])})
	}
	i = i + 1 + escapes
	return
}

// state after key=
func stateValue(p *Path, data string) (f stateFunc, i int) {
	var fl float64
	var escapes int

	r := strings.NewReader(data)
	x, _ := fmt.Fscan(r, &fl)

	if x > 0 {
		i = len(data) - r.Len()
	}

L:
	for ; i < len(data); i++ {
		switch data[i] {
		case '\\':
			i++
			if i >= len(data) {
				f = stateEscape
				return
			}
			data = data[:i-1] + data[i:]
			i--
			escapes++
		case '.':
			f = stateParent
			break L
		}
	}

	appendPathSel(p, &pathVal{data[:i]})

	i = i + 1 + escapes

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

	sel, err := newPathIndex(data[:i])

	if err != nil {
		return addError(err.Error()), i
	}

	appendPathSel(p, sel)

	return stateArrayEnd, i + 1
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
	if len(data) > 1 && data[:2] == "**" {
		appendPathSel(p, &pathRec{})
		i = i + 2
	}

	if len(data) > i {
		if data[i] == '.' {
			f = stateSep
		} else {
			f = addError("expected seperator character")
		}
		i++
	}

	return
}

// state after ..
func stateParent(p *Path, data string) (f stateFunc, i int) {
	if !(len(data) > 0 && data[0] == '.') {
		f = addError(`expected ".."`)
		return
	}
	appendPathSel(p, &pathParent{})
	i++
	return stateKey, i
}

// addError returns a stateFunc which sets the error.
func addError(msg string) stateFunc {
	return func(p *Path, data string) (f stateFunc, i int) {
		p.err = &syntaxError{"invalid path at " + p.p[:len(p.p)-len(data)] + " " + msg}
		return
	}
}

var stateEscape = addError("invalid escape character")
