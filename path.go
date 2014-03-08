package goj

import (
	"log"
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
	var arr []interface{}
	filterPath(d.v, arr, NewPath(s))
}

func filterPath(v interface{}, arr []interface{}, path *Path) (interface{}, bool) {
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			arr = append(arr, key)
			if val, ok := filterPath(val, arr, path); !ok {
				delete(x, key)
			} else {
				x[key] = val
			}
		}
		if len(x) == 0 {
			return nil, false
		}
	case []interface{}:
		// TODO delete is not working
		for i := 0; i < len(x); i++ {
			arr = append(arr, i)
			if _, ok := filterPath(x[i], arr, path); !ok {
				// delete i
				x = append(x[:i], x[i+1:]...)
			}
		}
		if len(x) == 0 {
			return nil, false
		} else {
			return x, true
		}
	default:
		arr = append(arr, x)
		log.Print(arr)
		if arr[len(arr)-2] == "price" {
			return x, true
		}
		return nil, false
	}
	return v, true
}

func filterKey(k interface{}, path *Path) bool {
	switch x := k.(type) {
	case string:
		return x == "price"
	case int, float32:
	}
	return true
}
