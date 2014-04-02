package goj

// import "log"

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func (d *Decoder) FilterOn(s string) error {
	p, err := NewPath(s, d.v)

	if err != nil {
		return err
	}

	tree := NewTree(d.v)
	tree.Branches(func(b Branch) {
		filterVal(b, p)
	})
	d.v = cleanBuild(p.r)

	return nil
}

func filterVal(b Branch, p *Path) {
	filterBranch(b, p.sel, func(leaf *Leaf) {
		leaf.Traverse(func(lf *Leaf) {
			p.r = buildBranch(lf.GetBranch(), p.r)
		})
	})
}

func filterBranch(b Branch, sel []pathSel, cb func(*Leaf)) {
	var i, j int
	for ; i < len(sel) && j <= len(b); i, j = i+1, j+1 {
		x := sel[i]
		switch x.(type) {
		case *pathRec:
		case *pathParent:
			if j > 0 {
				j--
			}
			if b[j].kind == leafVal {
				j--
			}
			parent := b[j].Parent()
			if parent != nil {
				parent.Branches(func(br Branch) {
					filterBranch(br, sel[i+1:], cb)
				})
			}
			return
		default:
			if j >= len(b) || !x.Equal(b[j]) {
				return
			}
		}
	}
	if j >= len(b) {
		return
	}
	cb(b[j])
}
