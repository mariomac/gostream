package stream

import (
	"github.com/mariomac/gostream/item"
	"strings"
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

func TestToMap(t *testing.T) {
	processed := OfMap(map[string]int{
		"Barcelona": 1,
		"Madrid":    2,
		"Paris":     3,
	}).Map(func(p item.Pair[string, int]) item.Pair[string, int] {
		p.Key = strings.ToLower(p.Key)
		if strings.Contains(p.Key, "i") {
			p.Val++
		}
		return p
	})

	assert.Equal(t, map[string]int{
		"barcelona": 1,
		"madrid":    3,
		"paris":     4,
	}, ToMap(processed))
}
