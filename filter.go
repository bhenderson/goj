package goj

import "log"

func (d *Decoder) FilterOn(s string) error {
	var arr []interface{}
	p, err := NewPath(s)

	if err != nil {
		return err
	}

	filterPath(d.v, arr, p.sel)

	return nil
}

func filterPath(v interface{}, arr, selector []interface{}) (interface{}, bool) {
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			arr = append(arr, key)
			log.Println("aaaaaaaa", arr, selector)
			if val, ok := filterPath(val, arr, selector); !ok {
				delete(x, key)
			} else {
				x[key] = val
			}
			arr = arr[1:]
		}
		if len(x) == 0 {
			return nil, false
		}
	case []interface{}:
		for i := 0; i < len(x); i++ {
			arr = append(arr, i)
			if _, ok := filterPath(x[i], arr, selector); !ok {
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
		if !filterVal(arr, selector) {
			return nil, false
		}
	}
	return v, true
}

func filterVal(arr, selector []interface{}) bool {
	return arr[len(arr)-2] == "price"
}
