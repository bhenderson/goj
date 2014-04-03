package goj

// import "log"

func buildBranch(b Branch, v interface{}) interface{} {
	if len(b) < 1 {
		return v
	}

	leaf := b[0]
	switch leaf.kind {
	case leafKey:
		if v == nil {
			v = make(map[string]interface{})
		}
		r := v.(map[string]interface{})
		key := leaf.val.(string)
		val := r[key]
		r[key] = buildBranch(b[1:], val)
		v = r
	case leafIdx:
		if v == nil {
			v = make(map[int]interface{})
		}
		r := v.(map[int]interface{})
		key := leaf.val.(int)
		val := r[key]
		r[key] = buildBranch(b[1:], val)
		v = r
	case leafVal:
		v = leaf.val
	}
	return v
}

func cleanBuild(v interface{}) interface{} {
	switch x := v.(type) {
	case map[int]interface{}:
		i := 0
		r := make([]interface{}, len(x))
		for _, val := range x {
			r[i] = cleanBuild(val)
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
