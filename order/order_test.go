package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnoreCase(t *testing.T) {
	assert.Negative(t, IgnoreCase("aaa", "Bbb"))
	assert.Positive(t, IgnoreCase("Hello", "amici"))
	assert.Zero(t, IgnoreCase("LoloLo", "lOlolo"))
}

func TestNatural(t *testing.T) {
	assert.Negative(t, Natural(1, 2))
	assert.Zero(t, Natural("aaa", "aaa"))
	assert.Positive(t, Natural(3.2, 1.7))
}

func TestInverse(t *testing.T) {
	assert.Positive(t, Inverse(Natural[int])(1, 2))
	assert.Zero(t, Inverse(Natural[string])("aaa", "aaa"))
	assert.Negative(t, Inverse(Natural[float32])(3.2, 1.7))
}
