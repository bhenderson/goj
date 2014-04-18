package goj

import "fmt"

type leafFunc func(n *Leaf)

// TODO not sure what I want the API to do with this.
type branch []*Leaf

const trunkStr = "trunk"

// TODO the public API here needs cleanup.

// Leaf is the main data type to turn a json decoded object into a linked list.
// Tree is used just as a word to describe the object. Branches are a slice of
// Leaf pointers. A Branch itself is not a linked list, but each leaf links to
// it's parent. The value of a leaf is the key, index or value. The child of a
// leaf is the key or index value.
//
// 	obj := // json decoded object
// 	       // '{"a":{"b":"c"}}'
// 	Leaf(trunk, val: nil, child: obj)
// 	\_
// 	  Leaf(parent: trunk, val: "a", child: {"b": "c"})
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

// Parent returns a pointer to this leaf's parent leaf.
func (l *Leaf) Parent() *Leaf { return l.parent }

// GetBranch returns all the parents recursively of this leaf including this leaf.
// There is a special leaf type (the trunk) which will not be included in the
// branch, but it will still be referenced by the first leaf's parent. This is
// so that other branches can find each other.
func (l *Leaf) GetBranch() branch {
	b := make(branch, l.GetBranchLen())
	getBranch(l, b)
	return b
}

func (l *Leaf) GetBranchLen() (i int) {
	for p := l.Parent(); p != nil; p = p.Parent() {
		i++
	}
	return
}

// Traverse does a depth first search starting from leaf and yields the end
// node to the call back function.
func (l *Leaf) Traverse(cb leafFunc) { traverse(l, cb) }

// Child returns the interface{} that this leaf points to.
func (l *Leaf) Child() interface{} { return l.child }

// String is a convenience method for pretty printing a leaf.
func (l *Leaf) String() string {
	if l.kind == leafTrk {
		// trunk
		return trunkStr
	}
	return fmt.Sprint(l.val)
}

// all branches downstream of this leaf
func (l *Leaf) Branches(cb func(b branch)) {
	i := l.GetBranchLen()
	l.Traverse(func(leaf *Leaf) {
		cb(leaf.GetBranch()[i:])
	})
}

// PruneBranches takes a Path pointer and returns a pruned copy of this leaf's
// child.
func (l *Leaf) PruneBranches(p *Path) interface{} {
	v := copyZero(l.Child())

	cb := func(l2 *Leaf) {
		l2.Traverse(func(l3 *Leaf) {
			v = buildBranch(l3.GetBranch(), v)
		})
	}

	filterBranches(l, p.sel, cb)
	cleanBuild(v)
	return v
}

// NewTree takes an interface and returns a special Leaf pointer, the trunk,
// which references all other leaves.
func NewTree(v interface{}) *Leaf {
	return &Leaf{child: v, kind: leafTrk}
}

func getBranch(n *Leaf, b branch) {
	if len(b) < 1 {
		return
	}
	i := len(b) - 1
	b[i] = n
	if p := n.Parent(); p != nil {
		getBranch(p, b[:i])
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
