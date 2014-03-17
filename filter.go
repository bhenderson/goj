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

	filterPath(d.v, arr, p.sel)

	return nil
}

func filterPath(v interface{}, arr, selector []pathSel) {
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			wrap("key  ", &arr, pathKey{key}, func() {
				filterPath(val, arr, selector)
			})
		}
	case []interface{}:
		for i := 0; i < len(x); i++ {
			wrap("index", &arr, pathIndex{i}, func() {
				filterPath(x[i], arr, selector)
			})
		}
	default:
		wrap("value", &arr, pathVal{x}, func() {
			filterVal(arr, selector)
		})
	}
}

func wrap(msg string, arr *[]pathSel, v pathSel, cb func()) {
	pushState(arr, v)
	cb()
	log.Println(msg, arr)
	popState(arr)
}

func pushState(arr *[]pathSel, v pathSel) {
	*arr = append(*arr, v)
}

func popState(arr *[]pathSel) {
	*arr = (*arr)[:len(*arr)-1]
}

func filterVal(arr, selector []pathSel) bool {
	return false // arr[len(arr)-2] == "price"
}
