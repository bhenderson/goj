package goj

import "log"

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func (d *Decoder) FilterOn(s string) error {
	// log.SetOutput(new(NullWriter))

	var arr []pathSel
	p, err := NewPath(s, d.v)

	if err != nil {
		return err
	}

	filterPath(d.v, &arr, p)

	return nil
}

func filterPath(v interface{}, arrPtr *[]pathSel, p *Path) {
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			wrap("key  ", arrPtr, pathKey{key}, func() {
				filterPath(val, arrPtr, p)
			})
		}
	case []interface{}:
		for i := 0; i < len(x); i++ {
			wrap("index", arrPtr, pathIndex{i}, func() {
				filterPath(x[i], arrPtr, p)
			})
		}
	default:
		wrap("value", arrPtr, pathVal{x}, func() {
			filterVal(arrPtr, p)
		})
	}
}

func wrap(msg string, arrPtr *[]pathSel, v pathSel, cb func()) {
	pushState(arrPtr, v)
	cb()
	// log.Println(msg, arrPtr)
	popState(arrPtr)
}

func pushState(arrPtr *[]pathSel, v pathSel) {
	*arrPtr = append(*arrPtr, v)
}

func popState(arrPtr *[]pathSel) {
	*arrPtr = (*arrPtr)[:len(*arrPtr)-1]
}

func filterVal(arrPtr *[]pathSel, p *Path) {
	arr := *arrPtr
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

	log.Println("aaaaa", arr)
	v := findPath(&arr, p.v)
	log.Println("bbbbb", v)

	if i < len(p.sel) {
		p2 := &Path{sel: p.sel[i:], v: v}
		var arr2 []pathSel
		filterPath(v, &arr2, p2)
		log.Println("ddddd", arr)
	} else {
		arr = append(arr, pathVal{v})
		log.Println("ccccc", arr)
	}
}

func findPath(arrPtr *[]pathSel, v interface{}) interface{} {
	arr := *arrPtr
	for _, sel := range arr {
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
