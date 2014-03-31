package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var input = `{
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

func TestPath_FilterOn(t *testing.T) {
	var exp, m string
	var e, a interface{}

	// parent
	exp = `{
		"store": {
			"bicycles": [
				{
					"color": "red"
				}
			],
			"truck": {
				"color": "yellow"
			}
		}
	}`
	e, a, m = testFilterOn(t, exp, input, "**.price=3.99..color")
	assert.Equal(t, e, a, m)

	// wrong filter, returns nil
	exp = `{}`
	e, a, m = testFilterOn(t, exp, input, "blah")
	assert.Equal(t, e, a, m)

	// recursive without value
	exp = `{
		"store": {
			"bicycles": [
				{
					"color": "red"
				},
				{
					"color": "blue"
				}
			],
			"truck": {
				"color": "yellow"
			}
		}
	}`
	e, a, m = testFilterOn(t, exp, input, "store.**.color")
	assert.Equal(t, e, a, m)

	// top level key
	e, a, m = testFilterOn(t, input, input, "store")
	assert.Equal(t, e, a, m)

	// end with recursive
	exp = `{}`
	e, a, m = testFilterOn(t, exp, input, "store.**")
	assert.Equal(t, e, a, m)

	// recursive with missing path
	exp = `{}`
	e, a, m = testFilterOn(t, exp, input, "store.**.blah")
	assert.Equal(t, e, a, m)

	// index
	exp = `{
		"store": {
			"bicycles": [
				{
					"color": "blue",
					"price": 2.99
				}
			]
		}
	}`
	e, a, m = testFilterOn(t, exp, input, "store.bicycles[1]")
	assert.Equal(t, e, a, m)

	exp = `{
		"store": {
			"bicycles": [
				{
					"price": 3.99
				}
			],
			"truck": {
				"price": 3.99
			}
		}
	}`
	e, a, m = testFilterOn(t, exp, input, "**.=3.99")
	assert.Equal(t, e, a, m)

	exp = `{
		"store": {
			"bicycles": [
				{
					"color": "red",
					"price": 3.99
				}
			]
		}
	}`
	e, a, m = testFilterOn(t, exp, input, "store.bicycles[0].price=3.99..")
	assert.Equal(t, e, a, m)

	exp = `{
		"store": {
			"bicycles": [
				{
					"color": "red"
				}
			]
		}
	}`
	e, a, m = testFilterOn(t, exp, input, "store.bicycles[0].price=3.99..color")
	assert.Equal(t, e, a, m)

	e, a, m = testFilterOn(t, `{}`, input, "store.bicycles..blah")
	assert.Equal(t, e, a, m)

	exp = `{
		"store": {
			"bicycles": [
				{
					"color": "blue"
				}
			]
		}
	}`
	e, a, m = testFilterOn(t, exp, input, "store.bicycles[0].color....[1].color")
	assert.Equal(t, e, a, m)

	e, a, m = testFilterOn(t, `{}`, input, "store.bicycles[0].color....[1].blah")
	assert.Equal(t, e, a, m)
}

func Test_filterRec(t *testing.T) {
	// arr := []pathSel{
	// &pathKey{"store"},
	// &pathKey{"bicycles"},
	// &pathIdx{0, 1},
	// &pathKey{"price"},
	// &pathVal{3.99},
	// }

	// // **.=3.99..color
	// sel := []pathSel{
	// &pathRec{},
	// pathStar,
	// &pathVal{"3.99"},
	// &pathParent{},
	// &pathKey{"color"},
	// }

	// assert.True(t, filterRec(arr, sel[1:]))
}

func testFilterOn(t *testing.T, exp, input string, filter string) (e, a interface{}, m string) {
	d1 := testDecoder(t, exp)
	d2 := testDecoder(t, input)
	err := d2.FilterOn(filter)

	if err != nil {
		t.Fatal(err)
	}

	// t.Log(d1, d2)
	return d1.v, d2.v, filter
}
