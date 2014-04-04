package goj

import "fmt"

type leafFunc func(n *Leaf)

// TODO not sure what I want the API to do with this.
type branch []*Leaf

const trunkStr = "trunk"

type Leaf struct {
	parent *Leaf
	child  interface{}
	val    interface{}
	max    int
	kind   int
}

const (
	leafTrk = iota
	leafKey
	leafIdx
	leafVal
)

func (l *Leaf) Parent() *Leaf         { return l.parent }
func (l *Leaf) GetBranch() (b branch) { return getBranch(l) }

// Traverse does a depth first search starting from leaf and yields the end
// node to the call back function.
func (l *Leaf) Traverse(cb leafFunc) { traverse(l, cb) }
func (l *Leaf) Child() interface{}   { return l.child }
func (l *Leaf) String() string {
	if l.kind == leafTrk {
		// trunk
		return trunkStr
	}
	return fmt.Sprint(l.val)
}

// all branches downstream of this leaf
func (l *Leaf) Branches(cb func(b branch)) {
	i := len(l.GetBranch())
	l.Traverse(func(leaf *Leaf) {
		cb(leaf.GetBranch()[i:])
	})
}

func (leaf *Leaf) PruneBranches(p *Path) interface{} {
	v := copyZero(leaf.Child())

	cb := func(l2 *Leaf) {
		l2.Traverse(func(l3 *Leaf) {
			v = buildBranch(l3.GetBranch(), v)
		})
	}

	filterBranches(leaf, p.sel, cb)
	cleanBuild(v)
	return v
}

func NewTree(v interface{}) *Leaf {
	return &Leaf{child: v, kind: leafTrk}
}

func getBranch(n *Leaf) (b branch) {
	if p := n.Parent(); p != nil {
		b = p.GetBranch()
		b = append(b, n)
		return b
	} else if n.kind == leafTrk {
		return branch{}
	} else {
		return branch{n}
	}
}

func traverse(parent *Leaf, cb leafFunc) {
	var leaf *Leaf
	v := parent.Child()

	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			leaf = &Leaf{
				parent: parent,
				child:  val,
				val:    key,
				kind:   leafKey,
			}
			traverse(leaf, cb)
		}
	case []interface{}:
		l := len(x)
		for i, val := range x {
			leaf = &Leaf{
				parent: parent,
				child:  val,
				val:    i,
				max:    l,
				kind:   leafIdx,
			}
			traverse(leaf, cb)
		}
	default:
		leaf = &Leaf{
			parent: parent,
			val:    x,
			kind:   leafVal,
		}
		cb(leaf)
	}
}

func copyZero(v interface{}) interface{} {
	switch v.(type) {
	case map[string]interface{}:
		return map[string]interface{}{}
	case []interface{}:
		return []interface{}{}
	}
	return nil
}
