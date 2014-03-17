package goj

import "log"

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func (d *Decoder) FilterOn(s string) error {
	// log.SetOutput(new(NullWriter))

	var arr []pathSel
	p, err := NewPath(s)

	if err != nil {
		return err
	}

	filterPath(d.v, arr, &p.sel)

	return nil
}

func filterPath(v interface{}, arr []pathSel, sel *[]pathSel) {
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			wrap("key  ", &arr, pathKey{key}, func() {
				filterPath(val, arr, sel)
			})
		}
	case []interface{}:
		for i := 0; i < len(x); i++ {
			wrap("index", &arr, pathIndex{i}, func() {
				filterPath(x[i], arr, sel)
			})
		}
	default:
		wrap("value", &arr, pathVal{x}, func() {
			filterVal(arr, sel)
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

func filterVal(arr []pathSel, sel *[]pathSel) {
	var i, j int
	var x, y pathSel
	for ; i < len(*sel) && j < len(arr); i, j = i+1, j+1 {
		x = (*sel)[i]
		y = arr[j]

		switch x.(type) {
		case pathRec:
			i++
			x = (*sel)[i]
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
	if i < len(*sel) {
		x = (*sel)[i]
		if _, ok := x.(pathParent); ok {
			j--
			if _, ok = (arr[j]).(pathVal); ok {
				j--
			}
		}
	}

	log.Println("aaaaa", j, arr[:j])
}
