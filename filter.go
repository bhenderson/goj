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

	tree := NewTree(d.v)
	tree.Branches(func(b Branch) {
		filterVal(b, p)
	})
	d.v = cleanBuild(p.r)

	return nil
}

func filterVal(b Branch, p *Path) {
	leaf := filterBranch(b, p.sel, p)
	if leaf != nil {
		p.r = buildBranch(leaf.GetBranch(), p.r)
	}
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

func filterBranch(b Branch, sel []pathSel, p *Path) *Leaf {
	var i, j int
	for ; i < len(sel) && j <= len(b); i, j = i+1, j+1 {
		x := sel[i]
		switch x.(type) {
		case *pathRec:
		case *pathParent:
			// b = filterParent(b[:j], sel[i+1:], p)
			// return b[j]
		default:
			if j >= len(b) || !x.Equal(b[j]) {
				return nil
			}
		}
	}
	if j == len(b) {
		return nil
	}
	return b[j]
}

func filterParent(b Branch, sel []pathSel, p *Path) Branch {
	/* j := len(b) - 1
	 * if b[j].kind == leafVal {
	 *     j--
	 * }
	 * b = b[:j]
	 * v := b[j-1].val
	 * if len(sel) > 0 {
	 *     sel2 := make([]pathSel, j)
	 *     for i := range sel2 {
	 *         sel2[i] = pathStar
	 *     }
	 *     sel2 = append(sel2, sel...)
	 *     p2 := &Path{sel: sel2, v: p.v}
	 *     // filterPath(p.v, []pathSel{}, p2)
	 *     v = cleanBuild(p2.r)
	 *     v = findPath(b, v)
	 *     if v == nil {
	 *         return []pathSel{}
	 *     }
	 *     log.Printf("%V", v)
	 * }
	 * b = append(b, &pathVal{v}) */
	return b
}

func filterRec(arr, sel []pathSel) bool {
	// last element of arr is pathVal
	/*     var j int
	 *     for j, _ = range sel {
	 *         if _, ok := sel[j].(*pathVal); ok {
	 *             break
	 *         }
	 *     }
	 *     sel = sel[:j+1]
	 *
	 *     for i, y := range arr {
	 *         sel[i].Equal(y)
	 *     } */

	return true
}
