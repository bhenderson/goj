package goj

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
