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
	Traverse(d.v, nil, func(n Noder) {
		i++
		t.Log(n.GetBranch())
	})

	assert.Equal(t, 6, i)
}
