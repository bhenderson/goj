package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPath_FilterOn(t *testing.T) {
	input := `{
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

	exp := `{
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

	d1 := testDecoder(t, exp)
	d2 := testDecoder(t, input)
	err := d2.FilterOn("**.price=3.99..color")

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

	d1 = testDecoder(t, input)
	d2 = testDecoder(t, input)
	err = d2.FilterOn("store")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(d1, d2)
	assert.Equal(t, d1.v, d2.v)
}
