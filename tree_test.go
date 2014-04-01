package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTraverse(t *testing.T) {
	var input = `{
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

	d := testDecoder(t, input)

	tree := NewTree(d.v)

	var leaf Leaf

	var i int
	tree.Traverse(func(n Leaf) {
		i++
		t.Log(n.GetBranch())
		leaf = n
	})

	t.Log("break")

	leaf = leaf.Parent()
	leaf = leaf.Parent()
	leaf = leaf.Parent()

	leaf.Traverse(func(n Leaf) {
		i++
		t.Log(n.GetBranch())
		leaf = n
	})

	assert.Equal(t, 14, i)
}
