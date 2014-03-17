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
		wrap("value", &arr, pathVal{x}, func() {
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
	for ; i < len(p.sel) && j < len(arr); i++ {
		// **
		x = p.sel[i]
		y = arr[j]
		// if pathRec, skip next one
		if _, ok := x.(pathRec); ok {
			i++
			x = p.sel[i]
			// price
			if !x.Equal(y) {
				i = i - 2 // retry
			}
		} else if !x.Equal(y) {
			return
		}
		j++
	}
	if !x.Equal(y) {
		return
	}
	log.Println("aaaaa", arr)
}
