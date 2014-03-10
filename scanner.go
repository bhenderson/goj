package goj

import (
	"strconv"
)

type stateFunc func(*Path, string) (stateFunc, int)

func (p *Path) parse() {
	f := stateKey
	data := p.p
	var i, x int

	for {
		if f == nil || i > len(data) {
			break
		}
		f, x = f(p, data[i:])
		i = i + x
	}
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

	p.sel[len(p.sel)-1] = pair
	return
}

// state after key[
func stateArray(p *Path, data string) (f stateFunc, i int) {
L1:
	for ; i < len(data); i++ {
		switch data[i] {
		case ']':
			break L1
		}
	}

	x, _ := strconv.Atoi(data[:i-1])
	p.sel = append(p.sel, x)
	i++

L2:
	for ; i < len(data); i++ {
		switch data[i] {
		case '\\':
			i++
		case '.':
			f = stateSep
			break L2
			// default error
		}
	}

	return stateSep, i + 1
}

// state after key
func stateSep(p *Path, data string) (f stateFunc, i int) {
	f = stateKey
	if len(data) > 0 {
		if data[0] == '.' {
			f = stateParent
		}
	}
	return
}

// look for **
func stateRecursive(p *Path, data string) (f stateFunc, i int) {
	return
}

func stateParent(p *Path, data string) (f stateFunc, i int) {
	p.sel = append(p.sel, "..")
	return stateKey, 1
}
