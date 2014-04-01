package goj

import "fmt"

type leafFunc func(n Leaf)

type Leaf interface {
	Child() interface{}
	GetBranch() Branch
	Parent() Leaf
	Traverse(leafFunc)
}

type Branch []Leaf

const trunkStr = "trunk"

// Trunk
type Trunk struct {
	child interface{}
}

func (n *Trunk) Parent() Leaf          { return nil }
func (n *Trunk) GetBranch() (b Branch) { return Branch{} }
func (n *Trunk) Traverse(cb leafFunc)  { traverse(n, cb) }
func (n *Trunk) Child() interface{}    { return n.child }
func (n *Trunk) String() string        { return trunkStr }

// LeafKey
type LeafKey struct {
	parent Leaf
	child  interface{}
	val    string
}

func (n *LeafKey) Parent() Leaf          { return n.parent }
func (n *LeafKey) GetBranch() (b Branch) { return getBranch(n) }
func (n *LeafKey) Traverse(cb leafFunc)  { traverse(n, cb) }
func (n *LeafKey) Child() interface{}    { return n.child }
func (n *LeafKey) String() string        { return n.val }

// LeafIdx
type LeafIdx struct {
	parent   Leaf
	child    interface{}
	val, max int
}

func (n *LeafIdx) Parent() Leaf          { return n.parent }
func (n *LeafIdx) GetBranch() (b Branch) { return getBranch(n) }
func (n *LeafIdx) Traverse(cb leafFunc)  { traverse(n, cb) }
func (n *LeafIdx) Child() interface{}    { return n.child }
func (n *LeafIdx) String() string        { return fmt.Sprint(n.val) }

// LeafVal
type LeafVal struct {
	parent Leaf
	val    interface{}
}

func (n *LeafVal) Parent() Leaf          { return n.parent }
func (n *LeafVal) GetBranch() (b Branch) { return getBranch(n) }
func (n *LeafVal) Traverse(cb leafFunc)  { traverse(n, cb) }
func (n *LeafVal) Child() interface{}    { return nil }
func (n *LeafVal) String() string        { return fmt.Sprint(n.val) }

func NewTree(v interface{}) Leaf {
	return &Trunk{v}
}

func getBranch(n Leaf) (b Branch) {
	if p := n.Parent(); p != nil {
		b = p.GetBranch()
		b = append(b, n)
		return b
	} else {
		return Branch{n}
	}
}

func traverse(parent Leaf, cb leafFunc) {
	var leaf Leaf
	v := parent.Child()

	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			leaf = &LeafKey{
				parent: parent,
				child:  val,
				val:    key,
			}
			traverse(leaf, cb)
		}
	case []interface{}:
		l := len(x)
		for i, val := range x {
			leaf = &LeafIdx{
				parent: parent,
				child:  val,
				val:    i,
				max:    l,
			}
			traverse(leaf, cb)
		}
	default:
		leaf = &LeafVal{
			parent: parent,
			val:    x,
		}
		cb(leaf)
	}
}
