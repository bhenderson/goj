package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var tree = `{
	"store": {
		"bicycles": [
			{
				"color": "red",
				"price": 3.99
			},
			{
				"color": "blue",
				"price": 2.99
			}
		],
		"truck": {
			"color": "yellow",
			"price": 3.99
		}
	},
	"grocers": {
		"color": "black"
	}
}`

func TestLeaf_Traverse(t *testing.T) {
	d := testDecoder(t, tree)

	tree := NewTree(d.v)

	var leaf *Leaf

	var i int
	tree.Traverse(func(n *Leaf) {
		i++
		t.Log(n.GetBranch())
		leaf = n
	})

	assert.Equal(t, 3, len(leaf.GetBranch()))

	leaf = leaf.Parent()
	leaf = leaf.Parent()
	leaf = leaf.Parent()

	t.Log(leaf)

	leaf.Traverse(func(n *Leaf) {
		i++
		t.Log(n.GetBranch())
		leaf = n
	})

	assert.Equal(t, 14, i)
}

func TestLeaf_Branches(t *testing.T) {
	d := testDecoder(t, tree)
	tree := NewTree(d.v)

	var leaf *Leaf
	var i int
	tree.Traverse(func(l *Leaf) {
		if i == 2 {
			leaf = l
		}
		i++
	})

	leaf = leaf.Parent()
	leaf = leaf.Parent()

	i = 0
	leaf.Branches(func(b Branch) {
		i++
		t.Log(b)
	})

	assert.Equal(t, 2, i)
}
