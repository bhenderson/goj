package goj

import (
	"fmt"
)

type nodeFunc func(n Node)

type Node interface {
	Parent() Node
	GetBranch() Branch
	Traverse(nodeFunc)
	String() string
}

type Branch []Node

func getBranch(n Node) (b Branch) {
	if p := n.Parent(); p != nil {
		b = p.GetBranch()
		b = append(b, n)
		return b
	} else {
		return Branch{n}
	}
}

func NewTree(v interface{}) Node {
	return &Trunk{v}
}

// Trunk
type Trunk struct {
	child interface{}
}

func (n *Trunk) Parent() Node { return nil }

func (n *Trunk) GetBranch() (b Branch) {
	return Branch{}
}

func (n *Trunk) Traverse(cb nodeFunc) {
	Traverse(n.child, n, cb)
}

func (n *Trunk) String() string {
	return "trunk ->"
}

// NodeKey
type NodeKey struct {
	parent Node
	child  interface{}
	val    string
}

func (n *NodeKey) Parent() Node {
	return n.parent
}

func (n *NodeKey) GetBranch() (b Branch) {
	return getBranch(n)
}

func (n *NodeKey) Traverse(cb nodeFunc) {
	Traverse(n.child, n, cb)
}

func (n *NodeKey) String() string {
	return n.val
}

// NodeIdx
type NodeIdx struct {
	parent   Node
	child    interface{}
	val, max int
}

func (n *NodeIdx) Parent() Node {
	return n.parent
}

func (n *NodeIdx) GetBranch() (b Branch) {
	return getBranch(n)
}

func (n *NodeIdx) Traverse(cb nodeFunc) {
	Traverse(n.child, n, cb)
}

func (n *NodeIdx) String() string {
	return fmt.Sprint(n.val)
}

// NodeVal
type NodeVal struct {
	parent Node
	child  interface{}
	val    interface{}
}

func (n *NodeVal) Parent() Node {
	return n.parent
}

func (n *NodeVal) GetBranch() (b Branch) {
	return getBranch(n)
}

func (n *NodeVal) Traverse(cb nodeFunc) {
	Traverse(n.child, n, cb)
}

func (n *NodeVal) String() string {
	return fmt.Sprint(n.val)
}

func Traverse(v interface{}, parent Node, cb nodeFunc) {
	var node Node
	switch x := v.(type) {
	case map[string]interface{}:
		for key, val := range x {
			node = &NodeKey{
				parent: parent,
				child:  val,
				val:    key,
			}
			Traverse(val, node, cb)
		}
	case []interface{}:
		l := len(x)
		for i, val := range x {
			node = &NodeIdx{
				parent: parent,
				child:  val,
				val:    i,
				max:    l,
			}
			Traverse(val, node, cb)
		}
	default:
		node = &NodeVal{
			parent: parent,
			val:    x,
		}
		cb(node)
	}
}
