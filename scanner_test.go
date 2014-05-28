package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPath_Compile(t *testing.T) {
	var exp, act []pathSel
	var msg string

	exp, act, msg = helpPath(`store`, &pathKey{"store"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`store|foo`, &pathKey{"store|foo"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`store\\|foo`, &pathKey{`store\|foo`})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`.store`, &pathKey{"store"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`store.books`, &pathKey{"store"}, &pathKey{"books"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`a=b`, &pathKey{"a"}, &pathVal{"b"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`=`, &pathKey{"*"}, &pathVal{""})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`a=b..c`, &pathKey{"a"}, &pathVal{"b"}, &pathParent{}, &pathKey{"c"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`=3.99..c`, &pathKey{"*"}, &pathVal{"3.99"}, &pathParent{}, &pathKey{"c"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`=3.9*..c`, &pathKey{"*"}, &pathVal{"3.9*"}, &pathParent{}, &pathKey{"c"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`a..b`, &pathKey{"a"}, &pathParent{}, &pathKey{"b"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[0]`, &pathIndex{[]int{0}})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`books[0].price`, &pathKey{"books"}, &pathIndex{[]int{0}}, &pathKey{"price"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`books[0]..[1].price`, &pathKey{"books"}, &pathIndex{[]int{0}}, &pathParent{}, &pathIndex{[]int{1}}, &pathKey{"price"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`**`, &pathRec{})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`**.price`, &pathRec{}, &pathKey{"price"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`**..price`, &pathRec{}, &pathParent{}, &pathKey{"price"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`**.price=3.99`, &pathRec{}, &pathKey{"price"}, &pathVal{"3.99"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`**.price=3.99*`, &pathRec{}, &pathKey{"price"}, &pathVal{"3.99*"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[*]`, &pathSlice{0, -1, 1})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[0:1]`, &pathSlice{0, 1, 1})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[-1:]`, &pathSlice{-1, -1, 1})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[0:10:2]`, &pathSlice{0, 10, 2})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[1,2,-1]`, &pathIndex{[]int{1, 2, -1}})
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

	exp, act, msg = helpPathErr(`[0a]`, `invalid path at [0a invalid array index`)
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPathErr(`**b`, `invalid path at **b expected seperator character`)
	assert.Equal(t, exp, act, msg)

}

func TestPath_CompileEscapeChar(t *testing.T) {
	exp, act, msg := helpPath(`store\.books.hard`, &pathKey{"store.books"}, &pathKey{"hard"})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`stor\[e\][0]`, &pathKey{"stor[e]"}, &pathIndex{[]int{0}})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`=hi\.mom..`, &pathKey{"*"}, &pathVal{"hi.mom"}, &pathParent{})
	assert.Equal(t, exp, act, msg)
}

func BenchmarkNewPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPath(`stor\[e\][0]..books[1].foo.bar.baz.**`)
	}
}

func helpPathErr(s, exp string) (e, a, m string) {
	_, err := NewPath(s)
	if err == nil {
		return exp, "no error occured", s
	}

	return exp, err.Error(), s
}

func helpPath(s string, exp ...pathSel) ([]pathSel, []pathSel, string) {
	p, e := NewPath(s)
	if e != nil {
		panic(e)
	}
	return exp, p.sel, s
}
