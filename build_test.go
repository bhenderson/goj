package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_buildPath(t *testing.T) {
	exp := `{
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
			},
			"books": [
				[ 1, 2 ],
				[ 3, 4 ]
			]
		}
	}`

	d := testDecoder(t, exp)
	tree := NewTree(d.v)

	var v interface{}

	tree.Traverse(func(l *Leaf) {
		v = buildBranch(l.GetBranch(), v)
	})
	v = cleanBuild(v)
	assert.Equal(t, d.v, v)
}
