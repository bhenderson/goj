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
				"color": "red",
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
				"color": "red",
				"price": 3.99
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
}
