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

	cb := func(leaf *Leaf) {
		leaf.Traverse(func(l *Leaf) {
			p.r = buildBranch(l.GetBranch(), p.r)
		})
	}

	filterBranches(tree, p.sel, cb)

	d.v = cleanBuild(p.r)

	return nil
}

func filterBranches(leaf *Leaf, sel []pathSel, cb func(*Leaf)) {
	leaf.Branches(func(b Branch) {
		filterBranch(b, sel, cb)
	})
}

func filterBranch(b Branch, sel []pathSel, cb func(*Leaf)) {
	var i, j int
	for ; i < len(sel) && j <= len(b); i, j = i+1, j+1 {
		x := sel[i]
		switch x.(type) {
		case *pathRec:
			// TODO return early if found
			filterBranch(b[j:], sel[i+1:], cb)
			i--
		case *pathParent:
			// TODO why?
			if j == len(b) {
				j--
			}
			filterParent(b[j], sel[i+1:], cb)
			return
		default:
			if j >= len(b) || !x.Equal(b[j]) {
				return
			}
		}
	}
	// TODO clean this up.
	if j > len(b) || len(b) < 1 {
		return
	}
	if j == len(b) {
		j--
	}
	cb(b[j])
}

func filterParent(leaf *Leaf, sel []pathSel, cb func(*Leaf)) {
	leaf = leaf.Parent()
	if leaf.kind == leafVal {
		leaf = leaf.Parent()
	}
	leaf = leaf.Parent()
	filterBranches(leaf, sel, cb)
}
