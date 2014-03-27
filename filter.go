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
	matchedPath := pathMatch(arr, p.sel, p)
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

func pathMatch(arr, sel []pathSel, p *Path) []pathSel {
	for i, j := 0, 0; i < len(sel); i, j = i+1, j+1 {
		x := sel[i]
		switch x.(type) {
		case *pathRec:
		case *pathParent:
		default:
			if j >= len(arr) || !x.Equal(arr[j]) {
				return []pathSel{}
			}
		}
	}
	return arr
}
