package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByKey(t *testing.T) {
	assert.Negative(t, ByKey[int, int](Natural)(1, 10, 2, 2))
	assert.Zero(t, ByKey[int, int](Natural)(2, 10, 2, 2))
	assert.Positive(t, ByKey[int, int](Natural)(2, 2, 1, 10))
}

func TestByVal(t *testing.T) {
	assert.Positive(t, ByVal[int, int](Natural)(1, 10, 2, 2))
	assert.Zero(t, ByVal[int, int](Natural)(1, 2, 3, 2))
	assert.Negative(t, ByVal[int, int](Natural)(2, 2, 1, 10))
}
