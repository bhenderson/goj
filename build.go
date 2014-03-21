package goj

func buildPath(arr []pathSel, v interface{}) interface{} {
	if len(arr) < 1 {
		return v
	}

	switch sel := arr[0].(type) {
	case pathKey:
		if v == nil {
			v = make(map[string]interface{})
		}
		r := v.(map[string]interface{})
		key := sel.val
		val := r[key]
		r[key] = buildPath(arr[1:], val)
		v = r
	case pathIndex:
		if v == nil {
			v = make(map[int]interface{})
		}
		r := v.(map[int]interface{})
		key := sel.val
		val := r[key]
		r[key] = buildPath(arr[1:], val)
		v = r
	case pathVal:
		v = sel.val
	}
	return v
}

func cleanBuild(v interface{}) interface{} {
	switch x := v.(type) {
	case map[int]interface{}:
		i := 0
		r := make([]interface{}, len(x))
		for _, val := range x {
			r[i] = val
			i++
		}
		v = r
	case map[string]interface{}:
		for key, val := range x {
			x[key] = cleanBuild(val)
		}
		v = x
	default:
		v = x
	}
	return v
}
