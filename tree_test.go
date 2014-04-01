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

	var node Node
	node = &NodeKey{child: d.v}

	var i int
	node.Traverse(func(n Node) {
		i++
		t.Log(n.GetBranch())
		node = n
	})

	t.Log("break")

	node = node.Parent()
	node = node.Parent()
	node = node.Parent()

	node.Traverse(func(n Node) {
		i++
		t.Log(n.GetBranch())
		node = n
	})

	assert.Equal(t, 14, i)
}
