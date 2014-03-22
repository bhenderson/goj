package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_pathKey_Equal(t *testing.T) {
	var p1, p2 pathSel

	p1 = pathKey{"store"}
	p2 = pathKey{"store"}
	assert.True(t, p2.Equal(p1))

	assert.False(t, p2.Equal(pathVal{"store"}))

	p2 = pathKey{"st*"}
	assert.True(t, p2.Equal(p1))

	p2 = pathKey{"st?re"}
	assert.True(t, p2.Equal(p1))

	p2 = pathKey{"st[aeiou]re"}
	assert.True(t, p2.Equal(p1))

	p2 = pathKey{"st[^a-np-z]re"}
	assert.True(t, p2.Equal(p1))

	p2 = pathKey{"blah"}
	assert.False(t, p2.Equal(p1))
}

func Test_pathVal_Equal(t *testing.T) {
	var p1, p2 pathSel

	p1 = pathVal{3.99}

	p2 = pathVal{"3.99"}
	assert.True(t, p2.Equal(p1))

	assert.False(t, p2.Equal(pathKey{"3.99"}))

	p2 = pathVal{"3*"}
	assert.True(t, p2.Equal(p1))

	p2 = pathVal{"3.[0-9]?"}
	assert.True(t, p2.Equal(p1))

	p2 = pathVal{"4.99"}
	assert.False(t, p2.Equal(p1))
}

func Test_pathIndex_Equal(t *testing.T) {
	var p1, p2 pathSel

	p1 = pathIdx{1, 1}

	p2 = pathIndex{"1"}
	assert.True(t, p2.Equal(p1))

	assert.False(t, p2.Equal(pathIndex{"1"}))

	p2 = pathIndex{"0"}
	assert.False(t, p2.Equal(p1))

}

func Test_pathSlice_Equal(t *testing.T) {
	var p1, p2 pathSel

	p1 = pathIdx{2, 10}

	// [0:10]
	p2 = pathSlice{0, 10, 1}
	assert.True(t, p2.Equal(p1))
	assert.False(t, p2.Equal(pathIdx{11, 12}))

	assert.False(t, p2.Equal(pathIndex{"2"}))

	// [0:10:2]
	p2 = pathSlice{0, 10, 2}
	assert.True(t, p2.Equal(p1))
	assert.False(t, p2.Equal(pathIdx{3, 11}))

	// [1:11:2] -> 1,3,5...11
	p2 = pathSlice{1, 11, 2}
	assert.True(t, p2.Equal(pathIdx{3, 11}))
	assert.False(t, p2.Equal(p1))

	p2 = pathSlice{5, 25, 5}
	assert.True(t, p2.Equal(pathIdx{15, 20}))
	assert.False(t, p2.Equal(pathIdx{9, 10}))

	// [-1:] last
	p2 = pathSlice{-1, -1, 1}
	assert.True(t, p2.Equal(pathIdx{9, 10}))
	assert.False(t, p2.Equal(p1))

	// [:] all
	p2 = pathSlice{0, -1, 1}
	assert.True(t, p2.Equal(p1))
}

func Test_newPathSet(t *testing.T) {
	p := newPathSet("0,1,-1")
	assert.Equal(t, []int{0, 1, -1}, p.val)

	p = newPathSet("0")
	assert.Equal(t, []int{0}, p.val)
}

func Test_pathSet_Equal(t *testing.T) {
	var p1, p2 pathSel

	p1 = pathIdx{1, 10}

	p2 = newPathSet("0,1")
	assert.True(t, p2.Equal(p1))
	assert.True(t, p2.Equal(pathIdx{0, 10}))
	assert.False(t, p2.Equal(pathIdx{3, 10}))
}
