package goj

func filterOn(d *Decoder, s string) error {
	p, err := NewPath(s)

	if err != nil {
		return err
	}

	tree := NewTree(d.v)

	d.v = tree.PruneBranches(p.sel)

	return nil
}

// filterBranches yields the last leaf in a branch that matches sel.
func filterBranches(leaf *Leaf, sel []pathSel, cb func(*Leaf)) {
	leaf.Branches(func(b branch) {
		filterBranch(b, sel, cb)
	})
}

func filterBranch(b branch, sel []pathSel, cb func(*Leaf)) {
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
