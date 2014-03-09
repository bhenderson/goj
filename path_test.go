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

	exp, act, msg = helpPath(`[0]`, 0)
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`books[0]`, "books", 0)
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[0:1]`, PairSlice{0, 1})
	assert.Equal(t, exp, act, msg)

	exp, act, msg = helpPath(`[-1:]`, PairSlice{-1, nil})
	assert.Equal(t, exp, act, msg)

	// books[] <- what does that mean?

	exp, act, msg = helpPath(`[1,2,-1]`, []int{1, 2, -1})
	assert.Equal(t, exp, act, msg)

}

func helpPath(s string, exp ...interface{}) ([]interface{}, []interface{}, string) {
	p := NewPath(s)
	return exp, p.sel, s
}

func TestBlah(t *testing.T) {
	input := `{
		"store": {
			"bicycles": [
				{
					"color": "red",
					"price": 3.99
				},
				{
					"color": "blue"
				}
			]
		}
	}`

	exp := `{
		"store": {
			"bicycles": [
				{
					"price": 3.99
				}
			]
		}
	}`

	d1 := testDecoder(t, exp)
	d2 := testDecoder(t, input)
	d2.FilterOn("price")

	assert.Equal(t, d1.v, d2.v)
}
