package goj

import (
	"fmt"
)

type Noder interface {
	Parent() Noder
	GetBranch() Branch
	String() string
}

type Branch []Noder

func getBranch(n Noder) (b Branch) {
	if p := n.Parent(); p != nil {
		b = p.GetBranch()
		b = append(b, n)
		return b
	} else {
		return Branch{n}
	}
}

type NodeKey struct {
	parent Noder
	child  interface{}
	val    string
}

func (n *NodeKey) Parent() Noder {
	return n.parent
}

func (n *NodeKey) GetBranch() (b Branch) {
	return getBranch(n)
}

func (n *NodeKey) String() string {
	return n.val
}

type NodeIdx struct {
	parent   Noder
	child    interface{}
	val, max int
}

func (n *NodeIdx) Parent() Noder {
	return n.parent
}

func (n *NodeIdx) GetBranch() (b Branch) {
	return getBranch(n)
}

func (n *NodeIdx) String() string {
	return fmt.Sprint(n.val)
}

type NodeVal struct {
	parent Noder
	child  interface{}
	val    interface{}
}

func (n *NodeVal) Parent() Noder {
	return n.parent
}

func (n *NodeVal) GetBranch() (b Branch) {
	return getBranch(n)
}

func (n *NodeVal) String() string {
	return fmt.Sprint(n.val)
}

type nodeFunc func(n Noder)

func Traverse(v interface{}, parent Noder, cb nodeFunc) {
	var node Noder
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
