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
	var exp string
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
	e, a = testFilterOn(t, exp, input, "**.price=3.99..color")
	assert.Equal(t, e, a)

	// wrong filter, returns nil
	exp = `null`
	e, a = testFilterOn(t, exp, input, "blah")
	assert.Equal(t, e, a)

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
	e, a = testFilterOn(t, exp, input, "store.**.color")
	assert.Equal(t, e, a)

	// top level key
	e, a = testFilterOn(t, input, input, "store")
	assert.Equal(t, e, a)
}

func testFilterOn(t *testing.T, exp, input string, filter string) (e, a interface{}) {
	d1 := testDecoder(t, exp)
	d2 := testDecoder(t, input)
	err := d2.FilterOn(filter)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(d1, d2)
	return d1.v, d2.v
}
