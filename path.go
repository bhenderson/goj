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
	p.parse()
	return p
}

type Pair struct {
	key string
	val interface{}
}

type PairSlice struct {
	b, e, s interface{}
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
		for i := 0; i < len(x); i++ {
			arr = append(arr, i)
			if _, ok := filterPath(x[i], arr, path); !ok {
				// delete i
				x = append(x[:i], x[i+1:]...)
			}
		}
		if len(x) == 0 {
			return nil, false
		}
		v = x
	default:
		arr = append(arr, x)
		if !filterVal(arr, path) {
			return nil, false
		}
	}
	return v, true
}

func filterVal(arr []interface{}, path *Path) bool {
	return arr[len(arr)-2] == "price"
}
