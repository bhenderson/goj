package goj

import "fmt"

type leafFunc func(n *Leaf)

type Branch []*Leaf

const trunkStr = "trunk"

type Leaf struct {
	parent *Leaf
	child  interface{}
	val    interface{}
	max    int
}

func (n *Leaf) Parent() *Leaf         { return n.parent }
func (n *Leaf) GetBranch() (b Branch) { return getBranch(n) }
func (n *Leaf) Traverse(cb leafFunc)  { traverse(n, cb) }
func (n *Leaf) Child() interface{}    { return n.child }
func (n *Leaf) isTrunk() bool {
	return n.parent == nil && n.val == nil
}
func (n *Leaf) String() string {
	if n.isTrunk() {
		// trunk
		return trunkStr
	}
	return fmt.Sprint(n.val)
}

// all branches downstream of this leaf
func (l *Leaf) Branches(cb func(b Branch)) {
	i := len(l.GetBranch())
	l.Traverse(func(leaf *Leaf) {
		cb(leaf.GetBranch()[i:])
	})
}

func NewTree(v interface{}) *Leaf {
	return &Leaf{child: v}
}

func getBranch(n *Leaf) (b Branch) {
	if p := n.Parent(); p != nil {
		b = p.GetBranch()
		b = append(b, n)
		return b
	} else if n.isTrunk() {
		return Branch{}
	} else {
		return Branch{n}
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
			}
			traverse(leaf, cb)
		}
	default:
		leaf = &Leaf{
			parent: parent,
			val:    x,
		}
		cb(leaf)
	}
}
