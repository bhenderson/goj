package goj

import (
	"strconv"
)

type Path struct {
	p     string
	sel   []interface{}
	depth int
}

func NewPath(s string) *Path {
	p := &Path{p: s}
	p.compile()
	return p
}

func (p *Path) compile() {
	start := 0
	for i := 0; i < len(p.p); i++ {
		switch p.p[i] {
		case '=':
			m := Pair{}
			key := p.p[start:i]
			if len(key) == 0 {
				key = "*"
			}
			m.key = key
			val := p.p[i+1:]
			if len(val) == 0 {
				m.val = nil
			} else {
				m.val = val
			}
			p.sel = append(p.sel, m)
			start = len(p.p)
			break
		case '.', '[':
			if start != i {
				p.sel = append(p.sel, p.p[start:i])
			}
			start = i + 1
		case ']':
			p.sel = append(p.sel, col(p.p[start:i]))
			start = i + 1
		}
	}
	if start != len(p.p) {
		p.sel = append(p.sel, p.p[start:])
	}
}

type Pair struct {
	key string
	val interface{}
}

type PairSlice struct {
	b, e interface{}
}

func col(s string) interface{} {
	var start int
	var arr []int
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ':':
			var p PairSlice
			if i == 0 {
				p.b = nil
			} else {
				x, _ := strconv.Atoi(s[0:i])
				p.b = x
			}
			i++
			if i == len(s) {
				p.e = nil
			} else {
				x, _ := strconv.Atoi(s[i:])
				p.e = x
			}
			return p
		case ',':
			x, _ := strconv.Atoi(s[start:i])
			arr = append(arr, x)
			start = i + 1
		}
	}
	if start != 0 {
		x, _ := strconv.Atoi(s[start:])
		arr = append(arr, x)
	}
	if len(arr) != 0 {
		return arr
	}
	i, _ := strconv.Atoi(s)
	return i
}

func (d *Decoder) FilterOn(s string) {
	filterPath(d.v, NewPath(s))
}

func filterPath(v interface{}, path *Path) (suc bool) {
	path.depth++
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			if !filterKey(key, path) && filterPath(val, path) {
				delete(x, key)
			}
		}
		if len(x) == 0 {
			return false
		}
	case []interface{}:
		for i := 0; i < len(x); i++ {
			if !filterKey(i, path) && !filterPath(x[i], path) {
				// delete i
				x = append(x[:i], x[i+1:]...)
			}
		}
		if len(x) == 0 {
			return false
		}
	}
	return true
}

func filterKey(k interface{}, path *Path) bool {
	switch x := k.(type) {
	case string:
		return x == "price"
	case int, float32:
	}
	return true
}
