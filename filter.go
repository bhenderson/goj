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
			pushState(&arr, pathKey{key})
			filterPath(val, arr, selector)
			log.Println("key  ", arr)
			popState(&arr)
		}
	case []interface{}:
		for i := 0; i < len(x); i++ {
			pushState(&arr, pathIndex{i})
			filterPath(x[i], arr, selector)
			log.Println("index", arr)
			popState(&arr)
		}
	default:
		pushState(&arr, pathVal{x})
		filterVal(arr, selector)
		log.Println("value", arr)
		popState(&arr)
	}
}

func filterVal(arr, selector []pathSel) bool {
	return false // arr[len(arr)-2] == "price"
}

func pushState(arr *[]pathSel, v pathSel) {
	*arr = append(*arr, v)
}

func popState(arr *[]pathSel) {
	*arr = (*arr)[:len(*arr)-1]
}
