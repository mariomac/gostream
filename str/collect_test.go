package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
)

func TestAsSlice(t *testing.T) {
	// Testing AsSlice of concrete streams
	slice := OfSlice([]int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, slice.AsSlice())

	// Testing AsSlice of connectedStreams
}
