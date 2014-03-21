package goj

import (
	"fmt"
	// "log"
)

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func (d *Decoder) FilterOn(s string) error {
	// log.SetOutput(new(NullWriter))

	var arr []pathSel
	p, err := NewPath(s, d.v)

	if err != nil {
		return err
	}

	filterPath(d.v, arr, p)

	return nil
}

func filterPath(v interface{}, arr []pathSel, p *Path) {
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			wrap("key  ", &arr, pathKey{key}, func() {
				filterPath(val, arr, p)
			})
		}
	case []interface{}:
		for i := 0; i < len(x); i++ {
			wrap("index", &arr, pathIndex{i}, func() {
				filterPath(x[i], arr, p)
			})
		}
	default:
		wrap("value", &arr, pathVal{fmt.Sprint(x)}, func() {
			filterVal(arr, p)
		})
	}
}

func wrap(msg string, arr *[]pathSel, v pathSel, cb func()) {
	pushState(arr, v)
	cb()
	// log.Println(msg, arr)
	popState(arr)
}

func pushState(arr *[]pathSel, v pathSel) {
	*arr = append(*arr, v)
}

func popState(arr *[]pathSel) {
	*arr = (*arr)[:len(*arr)-1]
}

func filterVal(arr []pathSel, p *Path) {
	var i, j int
	var x, y pathSel

	for ; i < len(p.sel) && j < len(arr); i, j = i+1, j+1 {
		x = p.sel[i]
		y = arr[j]

		switch x.(type) {
		case pathRec:
			i++
			x = p.sel[i]
			if !x.Equal(y) {
				i = i - 2 // retry
			}
		default:
			if !x.Equal(y) {
				return
			}
		}
	}

	// was last compare true?
	if !x.Equal(y) {
		return
	}

	// eval parent
	if i < len(p.sel) {
		x = p.sel[i]
		if _, ok := x.(pathParent); ok {
			i, j = i+1, j-1
			if _, ok = (arr[j]).(pathVal); ok {
				j--
			}
		}
	}

	arr = arr[:j]
	v := findPath(&arr, p.v)

	if i >= len(p.sel) {
		p.res = arr
		if p.p == "" {
			return
		}
	} else {
		p2 := &Path{sel: p.sel[i:], v: v}
		var arr2 []pathSel
		filterPath(v, arr2, p2)
		arr = append(arr, p2.res...)
	}

	buildPath(arr, p)
}

func findPath(arrPtr *[]pathSel, v interface{}) interface{} {
	// TODO error checking
	for _, sel := range *arrPtr {
		switch x := v.(type) {
		case map[string]interface{}:
			v = x[(sel.(pathKey)).val]
		case []interface{}:
			v = x[(sel.(pathIndex)).val]
		default:
			v = x
		}
	}

	return v
}
