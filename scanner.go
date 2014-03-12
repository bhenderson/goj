package goj

import (
	"strconv"
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
func stateKey(p *Path, data string) (f stateFunc, i int) {
L:
	for ; i < len(data); i++ {
		switch data[i] {
		case '\\':
			// escape char. skip next value
			i++
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
		p.sel = append(p.sel, data[:i])
	}
	i++
	return
}

// state after key=
func stateValue(p *Path, data string) (f stateFunc, i int) {
	var pair Pair
	pair.key = p.sel[len(p.sel)-1].(string)

L:
	for ; i < len(data); i++ {
		switch data[i] {
		case '\\':
			i++
		case '.':
			f = stateParent
			break L
		}
	}

	if i != 0 {
		pair.val = data[:i]
	}

	i++

	p.sel[len(p.sel)-1] = pair
	return
}

// state after key[
func stateArray(p *Path, data string) (f stateFunc, i int) {
L:
	for ; i < len(data); i++ {
		switch data[i] {
		case '\\':
			i++
		case ']':
			break L
		}
	}

	x, e := strconv.Atoi(data[:i])
	if e != nil {
		return
	}
	p.sel = append(p.sel, x)
	i++

	if i < len(data) {
		if data[i] == '.' {
			f = stateSep
		} else {
			f = addError(`expected "."`)
		}
		i++
	}

	return f, i
}

// state after .
func stateSep(p *Path, data string) (f stateFunc, i int) {
	if len(data) > 0 && data[0] == '.' {
		return stateParent, i
	}
	return stateKey, i
}

// look for **
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
