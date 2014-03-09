package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDotSplit(t *testing.T) {
	exp, act := helpDotSplit(`.a.b\.c.d`, "a", `b\.c`, "d")
	assert.Equal(t, exp, act)

	exp, act = helpDotSplit(`a..b`, "a", "..", "b")
	assert.Equal(t, exp, act)
}

func helpDotSplit(s string, strs ...string) ([]string, []string) {
	return strs, dotSplit(s)
}
