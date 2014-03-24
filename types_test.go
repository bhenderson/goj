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

	p1 = pathIdx{1, 10}

	p2 = pathIndex{[]int{1}}
	assert.True(t, p2.Equal(p1))

	assert.False(t, p2.Equal(pathIndex{[]int{1}}))

	p2 = pathIndex{[]int{0}}
	assert.False(t, p2.Equal(p1))

	p2 = pathIndex{[]int{0, 1}}
	assert.True(t, p2.Equal(p1))
	assert.True(t, p2.Equal(pathIdx{0, 10}))
	assert.False(t, p2.Equal(pathIdx{3, 10}))
}

func Test_pathSlice_Equal(t *testing.T) {
	var p1, p2 pathSel

	p1 = pathIdx{2, 10}

	// [0:10]
	p2 = pathSlice{0, 10, 1}
	assert.True(t, p2.Equal(p1))
	assert.False(t, p2.Equal(pathIdx{11, 12}))

	assert.False(t, p2.Equal(pathIndex{[]int{2}}))

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

func Test_newPathIndex(t *testing.T) {
	sel := testNewPathIndex(t, "*")
	assert.Equal(t, pathSlice{0, -1, 1}, sel)

	sel = testNewPathIndex(t, ":")
	assert.Equal(t, pathSlice{0, -1, 1}, sel)

	sel = testNewPathIndex(t, "::")
	assert.Equal(t, pathSlice{0, -1, 1}, sel)

	sel = testNewPathIndex(t, "0:1")
	assert.Equal(t, pathSlice{0, 1, 1}, sel)

	sel = testNewPathIndex(t, "1:")
	assert.Equal(t, pathSlice{1, -1, 1}, sel)

	sel = testNewPathIndex(t, ":2")
	assert.Equal(t, pathSlice{0, 2, 1}, sel)

	sel = testNewPathIndex(t, "0:-1:2")
	assert.Equal(t, pathSlice{0, -1, 2}, sel)

	sel = testNewPathIndex(t, "::2")
	assert.Equal(t, pathSlice{0, -1, 2}, sel)

	sel = testNewPathIndex(t, "0")
	assert.Equal(t, pathIndex{[]int{0}}, sel)

	sel = testNewPathIndex(t, "0,1,-3")
	assert.Equal(t, pathIndex{[]int{0, 1, -3}}, sel)
}

func Test_newPathIndex_Errors(t *testing.T) {
	msg := testNewPathIndexError(t, "a")
	assert.Equal(t, "invalid array index", msg)

	msg = testNewPathIndexError(t, "0:a")
	assert.Equal(t, "invalid array index", msg)

	msg = testNewPathIndexError(t, "0:1:a")
	assert.Equal(t, "invalid array index", msg)

	msg = testNewPathIndexError(t, "0,a")
	assert.Equal(t, "invalid array index", msg)
}

func testNewPathIndex(t *testing.T, s string) pathSel {
	sel, err := newPathIndex(s)

	if err != nil {
		t.Fatal(s, err)
	}

	return sel
}

func testNewPathIndexError(t *testing.T, s string) string {
	sel, err := newPathIndex(s)
	if err == nil {
		t.Fatalf("%s %s %#v", s, "expected error, got", sel)
	}
	return err.Error()
}
