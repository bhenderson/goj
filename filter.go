package goj

// import "log"

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func (d *Decoder) FilterOn(s string) error {
	// log.SetOutput(new(NullWriter))

	var arr []pathSel
	p, err := NewPath(s, d.v)

	// log.Println(p)

	if err != nil {
		return err
	}

	filterPath(d.v, arr, p)
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

func filterCompare(arr, sels []pathSel) bool {
	if len(arr) > len(sels) {
		return false
	}

	for i, val := range arr {
		if !sels[i].Equal(val) {
			return false
		}
	}

	return true
}

func filterVal(arr []pathSel, p *Path) {
	var i, j int
	var x, y pathSel

	for ; i < len(p.sel) && j < len(arr); i, j = i+1, j+1 {
		x = p.sel[i]
		y = arr[j]

		switch x.(type) {
		case *pathRec:
			i++
			if i >= len(p.sel) {
				return
			}
			if !filterCompare(arr[j:], p.sel[i:]) {
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
	if i < len(p.sel) {
		x = p.sel[i]
		if _, ok := x.(*pathParent); ok {
			i, j = i+1, j-1
			if _, ok = (arr[j]).(*pathVal); ok {
				j--
			}
		}
	}

	arr = arr[:j]
	v := findPath(&arr, p.v)

	if i >= len(p.sel) {
		if len(arr) > 0 {
			last := arr[len(arr)-1]
			if _, ok := last.(*pathVal); !ok {
				arr = append(arr, &pathVal{v})
			}
		}
		if p.p == "" {
			p.res = arr

			return
		}
	} else {
		p2 := &Path{sel: p.sel[i:], v: v}
		var arr2 []pathSel
		filterPath(v, arr2, p2)
		arr = append(arr, p2.res...)
	}

	p.r = buildPath(arr, p.r)
}

func findPath(arrPtr *[]pathSel, v interface{}) interface{} {
	// TODO error checking
	for _, sel := range *arrPtr {
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
