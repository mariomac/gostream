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
