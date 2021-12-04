package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
)

func TestToSlice(t *testing.T) {
	// Testing ToSlice of concrete streams
	slice := OfSlice([]int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, slice.ToSlice())

	// Testing ToSlice of connectedStreams
	assert.Equal(t, []int{1, 2, 3},
		slice.Filter(func(n int) bool {
			return n <= 3
		}).ToSlice())
}
