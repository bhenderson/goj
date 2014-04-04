package goj

import "fmt"

type leafFunc func(n *Leaf)

// TODO should we remove this?
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

func (n *Leaf) Parent() *Leaf         { return n.parent }
func (n *Leaf) GetBranch() (b branch) { return getBranch(n) }
func (n *Leaf) Traverse(cb leafFunc)  { traverse(n, cb) }
func (n *Leaf) Child() interface{}    { return n.child }
func (n *Leaf) String() string {
	if n.kind == leafTrk {
		// trunk
		return trunkStr
	}
	return fmt.Sprint(n.val)
}

// all branches downstream of this leaf
func (l *Leaf) Branches(cb func(b branch)) {
	i := len(l.GetBranch())
	l.Traverse(func(leaf *Leaf) {
		cb(leaf.GetBranch()[i:])
	})
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
