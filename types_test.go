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

	p1 = pathIdx{1}

	p2 = pathIndex{"1"}
	assert.True(t, p2.Equal(p1))

	p2 = pathIndex{"0"}
	assert.False(t, p2.Equal(p1))
}
