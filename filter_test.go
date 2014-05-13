package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_filterOn(t *testing.T) {
	var exp, m string
	var e, a interface{}

	// parent
	// wrong filter, returns nil
	exp = `{}`
	e, a, m = testFilterOn(t, exp, input, "blah")
	assert.Equal(t, e, a, m)

	// top level key
	e, a, m = testFilterOn(t, input, input, "store")
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
			"truck": {
				"color": "yellow"
			}
		}
	}`
	e, a, m = testFilterOn(t, exp, input, "store.truck.color=yellow")
	assert.Equal(t, e, a, m)

	// end with recursive
	e, a, m = testFilterOn(t, input, input, "store.**")
	assert.Equal(t, e, a, m)

	// recursive with missing path
	exp = `{}`
	e, a, m = testFilterOn(t, exp, input, "store.**.blah")
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
					"color": "red"
				},
				{
					"color": "blue"
				}
			]
		}
	}`
	e, a, m = testFilterOn(t, exp, input, "store.bicycles[0].color....[:].color")
	assert.Equal(t, e, a, m)

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
	e, a, m = testFilterOn(t, exp, input, "**.bicycles[-1]")
	assert.Equal(t, e, a, m)

	e, a, m = testFilterOn(t, `{}`, input, "store.bicycles[0].color....[1].blah")
	assert.Equal(t, e, a, m)

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
	e, a, m = testFilterOn(t, exp, input, "**.price=2.99..")
	assert.Equal(t, e, a, m)

	input = `[1,2,3,{"a":"b"}]`
	exp = `[{"a":"b"}]`
	e, a, m = testFilterOn(t, exp, input, "**.a")
	assert.Equal(t, e, a, m)
}

func testFilterOn(t *testing.T, exp, input, filter string) (e, a interface{}, m string) {
	v := testVal(t, input)
	err := filterOn(v, filter)

	if err != nil {
		t.Fatal(err)
	}

	// t.Log(d1, d2)
	return testMarshal(t, exp), v.v, filter
}

func BenchmarkFilterOn(b *testing.B) {
	v := testVal(b, input)
	tree := NewTree(v.v)
	p, _ := NewPath("**.price..color")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.PruneBranches(p)
	}
}
