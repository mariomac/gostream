package stream

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeq2(t *testing.T) {
	input := OfMap(map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6"}).
		Filter(func(k int, v string) bool {
			return k%2 == 0
		})
	output := map[int]string{}
	for k, v := range input.Seq2() {
		output[k] = v
	}
	assert.Equal(t, map[int]string{2: "2", 4: "4", 6: "6"}, output)
}

func TestSeq2_MultipleInstances(t *testing.T) {
	str := OfMap(map[int]int{3: 1, 2: 2, 1: 3})
	var keys, vals []int
	for k, v := range str.Seq2() {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	slices.Sort(keys)
	slices.Sort(vals)
	assert.Equal(t, []int{1, 2, 3}, keys)
	assert.Equal(t, []int{1, 2, 3}, vals)
	// check that a new Iter instance returns a fresh iterator from the start
	var keys2, vals2 []int
	for k, v := range str.Seq2() {
		keys2 = append(keys2, k)
		vals2 = append(vals2, v)
	}
	slices.Sort(keys2)
	slices.Sort(vals2)
	assert.Equal(t, []int{1, 2, 3}, keys2)
	assert.Equal(t, []int{1, 2, 3}, vals2)
}
