package goj

// import "log"

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func (d *Decoder) FilterOn(s string) error {
	// log.SetOutput(new(NullWriter))

	p, err := NewPath(s, d.v)

	if err != nil {
		return err
	}

	// log.Printf("%V", p.sel)

	filterPath(d.v, []pathSel{}, p)
	d.v = cleanBuild(p.r)

	return nil
}

func filterPath(v interface{}, arr []pathSel, p *Path) {
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			wrap("key  ", &arr, &pathKey{key}, func() {
				filterPath(val, arr, p)
			})
		}
	case []interface{}:
		l := len(x)
		for i := 0; i < l; i++ {
			wrap("index", &arr, &pathIdx{i, l}, func() {
				filterPath(x[i], arr, p)
			})
		}
	default:
		wrap("value", &arr, &pathVal{x}, func() {
			filterVal(arr, p)
		})
	}
}

// TODO don't need pointer here I think.
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
	matchedPath := filterMatched(arr, p.sel, p)
	p.r = buildPath(matchedPath, p.r)
}

func findPath(arr []pathSel, v interface{}) interface{} {
	// TODO error checking
	for _, sel := range arr {
		switch x := v.(type) {
		case map[string]interface{}:
			v = x[(sel.(*pathKey)).val]
		case []interface{}:
			v = x[(sel.(*pathIdx)).val]
		default:
			v = x
		}
	}

	return v
}

func filterMatched(arr, sel []pathSel, p *Path) []pathSel {
	for i, j := 0, 0; i < len(sel) && j <= len(arr); i, j = i+1, j+1 {
		x := sel[i]
		switch x.(type) {
		case *pathRec:
		case *pathParent:
			arr = filterParent(arr[:j], sel[i+1:], p)
			return arr
		default:
			if j >= len(arr) || !x.Equal(arr[j]) {
				return []pathSel{}
			}
		}
	}
	return arr
}

func filterParent(arr, sel []pathSel, p *Path) []pathSel {
	j := len(arr) - 1
	if _, ok := arr[j].(*pathVal); ok {
		j--
	}
	arr = arr[:j]
	v := findPath(arr, p.v)
	if len(sel) > 0 {
		sel2 := make([]pathSel, j)
		for i := range sel2 {
			sel2[i] = pathStar
		}
		sel2 = append(sel2, sel...)
		p2 := &Path{sel: sel2, v: p.v}
		filterPath(p.v, []pathSel{}, p2)
		if p2.r == nil {
			return []pathSel{}
		}
		v = cleanBuild(p2.r)
		v = findPath(arr, v)
	}
	arr = append(arr, &pathVal{v})
	return arr
}

func filterRec(arr, sel []pathSel) bool {
	// last element of arr is pathVal
	var j int
	for j, _ = range sel {
		if _, ok := sel[j].(*pathVal); ok {
			break
		}
	}
	sel = sel[:j+1]

	for i, y := range arr {
		sel[i].Equal(y)
	}

	return true
}
