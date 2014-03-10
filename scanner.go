package goj

import (
	"strconv"
)

type stateFunc func(*Path, string)

func (p *Path) parse() {
	stateKey(p, p.p)
}

func stateKey(p *Path, data string) {
	var f stateFunc
	var i int
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
	if f != nil && i <= len(data) {
		f(p, data[i:])
	}
}

// state after key=
func stateValue(p *Path, data string) {
	var f stateFunc
	var pair Pair
	var i int
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

	p.sel[len(p.sel)-1] = pair

	if f != nil && i < len(data) {
		f(p, data[i:])
	}
}

// state after key[
func stateArray(p *Path, data string) {
	var i int

L:
	for ; i < len(data); i++ {
		switch data[i] {
		case ']':
			break L
		}
	}

	x, _ := strconv.Atoi(data[:i-1])
	p.sel = append(p.sel, x)

	stateSep(p, data[i+1:])
}

func stateSep(p *Path, data string) {
	if len(data) > 0 {
		if data[0] == '.' {
			stateParent(p, data[1:])
		} else {
			stateKey(p, data)
		}
	}
}

// look for **
func stateRecursive(p *Path, data string) {
}

func stateParent(p *Path, data string) {
	p.sel = append(p.sel, "..")
	stateKey(p, data[1:])
}
