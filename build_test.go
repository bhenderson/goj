package goj

import (
	"encoding/json"
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

	val := testVal(t, exp)
	tree := NewTree(val.v)

	var v interface{}

	tree.Traverse(func(l *Leaf) {
		v = buildBranch(l.GetBranch(), v)
	})
	v = cleanBuild(v)
	assert.Equal(t, val.v, v)
}

func Test_buildPath_array(t *testing.T) {
	exp := `[1,2,3]`

	var val interface{}
	json.Unmarshal([]byte(exp), &val)
	tree := NewTree(val)

	var v interface{}

	tree.Traverse(func(l *Leaf) {
		v = buildBranch(l.GetBranch(), v)
	})
	v = cleanBuild(v)
	assert.Equal(t, val, v)
}
