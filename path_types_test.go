package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_pathKey_Equal(t *testing.T) {
	var p pathSel
	var l *Leaf

	l = &Leaf{val: "blah", kind: leafKey}
	p = &pathKey{"store"}

	assert.False(t, p.Equal(l))

	l.val = "store"
	assert.True(t, p.Equal(l))

	p = &pathKey{"st*"}
	assert.True(t, p.Equal(l))

	p = &pathKey{"st?re"}
	assert.True(t, p.Equal(l))

	p = &pathKey{"st[aeiou]re"}
	assert.True(t, p.Equal(l))

	p = &pathKey{"st[^a-np-z]re"}
	assert.True(t, p.Equal(l))

	l = &Leaf{val: 0, max: 1, kind: leafIdx}
	assert.True(t, pathStar.Equal(l))
}

func Test_pathVal_Equal(t *testing.T) {
	var p pathSel
	var l *Leaf

	l = &Leaf{val: 3.99, kind: leafVal}
	p = &pathVal{"3.99"}

	assert.True(t, p.Equal(l))

	l.val = "3.99" // string
	assert.True(t, p.Equal(l))

	p = &pathVal{"3*"}
	assert.True(t, p.Equal(l))

	p = &pathVal{"3.[0-9]?"}
	assert.True(t, p.Equal(l))

	p = &pathVal{"4.99"}
	assert.False(t, p.Equal(l))

	l.val = nil
	p = &pathVal{""}
	assert.True(t, p.Equal(l))

	p = &pathVal{"null"}
	assert.False(t, p.Equal(l))

	p = &pathVal{"*"}
	assert.False(t, p.Equal(l))

	l.val = ""
	assert.True(t, p.Equal(l))
}

func Test_pathIndex_Equal(t *testing.T) {
	var p pathSel
	var l *Leaf

	l = &Leaf{val: 1, max: 10, kind: leafIdx}
	p = &pathIndex{[]int{1}}

	assert.True(t, p.Equal(l))

	l.kind = leafVal
	assert.False(t, p.Equal(l))
	l.kind = leafIdx

	p = &pathIndex{[]int{0}}
	assert.False(t, p.Equal(l))

	p = &pathIndex{[]int{0, 1}}
	assert.True(t, p.Equal(l))
	l.val = 0
	assert.True(t, p.Equal(l))
	l.val = 3
	assert.False(t, p.Equal(l))

	p = &pathIndex{[]int{-1}}
	l.val = 9
	assert.True(t, p.Equal(l))
}

func Test_pathSlice_Equal(t *testing.T) {
	var p pathSel
	var l *Leaf

	l = &Leaf{val: 2, max: 10, kind: leafIdx}
	// [0:10]
	p = &pathSlice{0, 10, 1}

	assert.True(t, p.Equal(l))
	l.kind = leafVal
	assert.False(t, p.Equal(l))
	l.kind = leafIdx

	// [0:10:2]
	p = &pathSlice{0, 10, 2}
	assert.True(t, p.Equal(l))
	l.val = 3
	assert.False(t, p.Equal(l))

	// [1:11:2] -> 1,3,5...11
	p = &pathSlice{1, 11, 2}
	l.val = 11
	assert.True(t, p.Equal(l))
	l.val = 10
	assert.False(t, p.Equal(l))

	p = &pathSlice{5, 25, 5}
	l.val, l.max = 15, 20
	assert.True(t, p.Equal(l))
	l.val, l.max = 9, 10
	assert.False(t, p.Equal(l))

	// [-1:] last
	p = &pathSlice{-1, -1, 1}
	assert.True(t, p.Equal(l))
	l.val = 8
	assert.False(t, p.Equal(l))

	// [:] all
	p = &pathSlice{0, -1, 1}
	assert.True(t, p.Equal(l))
}

func Test_newPathIndex(t *testing.T) {
	sel := testNewPathIndex(t, "*")
	assert.Equal(t, &pathSlice{0, -1, 1}, sel)

	sel = testNewPathIndex(t, ":")
	assert.Equal(t, &pathSlice{0, -1, 1}, sel)

	sel = testNewPathIndex(t, "::")
	assert.Equal(t, &pathSlice{0, -1, 1}, sel)

	sel = testNewPathIndex(t, "0:1")
	assert.Equal(t, &pathSlice{0, 1, 1}, sel)

	sel = testNewPathIndex(t, "1:")
	assert.Equal(t, &pathSlice{1, -1, 1}, sel)

	sel = testNewPathIndex(t, ":2")
	assert.Equal(t, &pathSlice{0, 2, 1}, sel)

	sel = testNewPathIndex(t, "0:-1:2")
	assert.Equal(t, &pathSlice{0, -1, 2}, sel)

	sel = testNewPathIndex(t, "::2")
	assert.Equal(t, &pathSlice{0, -1, 2}, sel)

	sel = testNewPathIndex(t, "0")
	assert.Equal(t, &pathIndex{[]int{0}}, sel)

	sel = testNewPathIndex(t, "0,1,-3")
	assert.Equal(t, &pathIndex{[]int{0, 1, -3}}, sel)
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

	msg = testNewPathIndexError(t, "0a1")
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
