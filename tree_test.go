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
		}
	}`

	d := testDecoder(t, input)

	var i int
	var node Node
	Traverse(d.v, nil, func(n Node) {
		i++
		t.Log(n.GetBranch())
		node = n
	})

	t.Log("break")

	node = node.Parent()
	node = node.Parent()

	node.Traverse(func(n Node) {
		i++
		t.Log(n.GetBranch())
		node = n
	})

	assert.Equal(t, 8, i)
}
