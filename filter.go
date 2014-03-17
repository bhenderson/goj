package goj

import "log"

func (d *Decoder) FilterOn(s string) error {
	var arr []pathSel
	p, err := NewPath(s)

	if err != nil {
		return err
	}

	filterPath(d.v, arr, p.sel)

	return nil
}

func filterPath(v interface{}, arr, selector []pathSel) (interface{}, bool) {
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			arr = append(arr, pathKey{key})
			if val, ok := filterPath(val, arr, selector); !ok {
				delete(x, key)
			} else {
				x[key] = val
			}
			log.Println("key  ", arr, selector)
			arr = arr[:len(arr)-1]
		}
		if len(x) == 0 {
			return nil, false
		}
	case []interface{}:
		for i := 0; i < len(x); i++ {
			arr = append(arr, pathIndex{i})
			if _, ok := filterPath(x[i], arr, selector); !ok {
				// delete i
				x = append(x[:i], x[i+1:]...)
			}
			log.Println("index", arr, selector)
			arr = arr[:len(arr)-1]
		}
		if len(x) == 0 {
			return nil, false
		}
		v = x
	default:
		arr = append(arr, pathVal{x})
		log.Println("value", arr, selector)
		unsetFilterValue(arr)
		if !filterVal(arr, selector) {
			return nil, false
		}
	}
	return v, true
}

func filterVal(arr, selector []pathSel) bool {
	// var i, j int
	// for ; i < len(arr); i++ {
	// sel := selector[j]
	// lhs, lok := sel.(string)
	// rhs, rok := arr[i].(string)
	// if lok && rok {
	// }
	// }
	// // log.Println("aaaaaaaa", arr, selector)
	return false // arr[len(arr)-2] == "price"
}

func setFilterValue(arr []pathSel, v interface{}) {
	arr = append(arr, pathVal{v})
}

func unsetFilterValue(arr []pathSel) {
	arr = arr[:len(arr)-1]
}
