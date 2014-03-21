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
	var err error
	var d1, d2 *Decoder

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

	// parent
	d1 = testDecoder(t, exp)
	d2 = testDecoder(t, input)
	err = d2.FilterOn("**.price=3.99..color")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, d1.v, d2.v)

	d2 = testDecoder(t, input)
	err = d2.FilterOn("blah")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, nil, d2.v)

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
	d1 = testDecoder(t, exp)
	d2 = testDecoder(t, input)
	err = d2.FilterOn("store.**.color")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(d1, d2)
	assert.Equal(t, d1.v, d2.v)

	// top level key
	d1 = testDecoder(t, input)
	d2 = testDecoder(t, input)
	err = d2.FilterOn("store")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(d1, d2)
	assert.Equal(t, d1.v, d2.v)
}
