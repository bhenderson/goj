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
			}
		}
	}`

	dec := testDecoder(t, exp)

	// index describes original index, not dest index.
	a1 := []pathSel{
		pathKey{"store"},
		pathKey{"bicycles"},
		pathIndex{1},
		pathKey{"color"},
		pathVal{"red"},
	}
	a2 := []pathSel{
		pathKey{"store"},
		pathKey{"bicycles"},
		pathIndex{1},
		pathKey{"price"},
		pathVal{3.99},
	}
	a3 := []pathSel{
		pathKey{"store"},
		pathKey{"bicycles"},
		pathIndex{3},
		pathKey{"color"},
		pathVal{"blue"},
	}
	a4 := []pathSel{
		pathKey{"store"},
		pathKey{"bicycles"},
		pathIndex{3},
		pathKey{"price"},
		pathVal{2.99},
	}
	a5 := []pathSel{
		pathKey{"store"},
		pathKey{"truck"},
		pathKey{"color"},
		pathVal{"yellow"},
	}
	a6 := []pathSel{
		pathKey{"store"},
		pathKey{"truck"},
		pathKey{"price"},
		pathVal{3.99},
	}

	var v interface{}
	v = buildPath(a1, v)
	v = buildPath(a2, v)
	v = buildPath(a3, v)
	v = buildPath(a4, v)
	v = buildPath(a5, v)
	v = buildPath(a6, v)
	v = cleanBuild(v)

	assert.Equal(t, dec.v, v)
}
