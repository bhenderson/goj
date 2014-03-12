package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPath_Compile(t *testing.T) {
	exp, act, msg := helpPath(`store`, "store")
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`.store`, "store")
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`store.books`, "store", "books")
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`a=b`, Pair{"a", "b"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`=`, Pair{"*", nil})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`a=b..c`, Pair{"a", "b"}, "..", "c")
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`=3.99..c`, Pair{"*", 3.99}, "..", "c")
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`=3.9*..c`, Pair{"*", "3.9*"}, "..", "c")
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`a..b`, "a", "..", "b")
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[0]`, 0)
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`books[0].price`, "books", 0, "price")
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`books[0]..[1].price`, "books", 0, "..", 1, "price")
	assert.Equal(t, exp, act, msg)

	// exp, act, msg = helpPath(`[*]`, 0)
	// assert.Equal(t, exp, act, msg)

	// exp, act, msg = helpPath(`[]`, 0)
	// assert.Equal(t, exp, act, msg)

	return

	exp, act, msg = helpPath(`[0:1]`, PairSlice{0, 1, 1})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[-1:]`, PairSlice{-1, nil, 1})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[0:10:2]`, PairSlice{0, 10, 2})
	assert.Equal(t, exp, act, msg)

	// books[] <- what does that mean?

	exp, act, msg = helpPath(`[1,2,-1]`, []int{1, 2, -1})
	assert.Equal(t, exp, act, msg)

}

func TestPath_CompileErrors(t *testing.T) {
	exp, act, msg := helpPathErr(`a=b.c`, `invalid path at a=b. expected ".."`)
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPathErr(`a=b\`, `invalid path at a=b\ invalid escape character`)
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPathErr(`[0]\`, `invalid path at [0]\ expected "."`)
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPathErr(`[0]a`, `invalid path at [0]a expected "."`)
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPathErr(`store\`, `invalid path at store\ invalid escape character`)
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPathErr(`[0a]`, `invalid path at [0a invalid index`)
	assert.Equal(t, exp, act, msg)

}

func TestPath_CompileEscapeChar(t *testing.T) {
	exp, act, msg := helpPath(`store\.books`, "store.books")
	assert.Equal(t, exp, act, msg)

}

func helpPathErr(s, exp string) (e, a, m string) {
	_, err := NewPath(s)
	if err == nil {
		return exp, "no error occured", s
	}

	return exp, err.Error(), s
}

func helpPath(s string, exp ...interface{}) ([]interface{}, []interface{}, string) {
	p, e := NewPath(s)
	if e != nil {
		panic(e.Error())
	}
	return exp, p.sel, s
}
