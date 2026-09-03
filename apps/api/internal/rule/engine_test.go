package rule

import (
	"testing"
	"math"
	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	assert.True(t, compare(10, 5, "gt"))
	assert.False(t, compare(3, 5, "gt"))
	assert.True(t, compare(5, 5, "gte"))
	assert.True(t, compare(3, 5, "lt"))
	assert.False(t, compare(7, 5, "lt"))
	assert.True(t, compare(5, 5, "eq"))
	assert.True(t, compare(5, 5.0001, "eq"))
	assert.True(t, compare(5, 7, "neq"))
	assert.False(t, compare(5, 5, "neq"))
	assert.True(t, compare(5, 3, "lte"))
	assert.False(t, compare(7, 5, "lte"))
	assert.True(t, compare(5.0, 5, "eq"))
	assert.InDelta(t, 0.0, math.Abs(5.0-5), 0.0001)
}
