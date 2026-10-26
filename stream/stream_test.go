package stream

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLazyOperation(t *testing.T) {
	var actions []string
	in := Of("hello", "my", "friend")
	filtered := in.Filter(func(s string) bool {
		actions = append(actions, "filter("+s+")")
		return strings.Contains(s, "e")
	})
	mapped := filtered.Map(func(s string) int {
		actions = append(actions, "map("+s+")")
		return len(s)
	})
	mapped.ForEach(func(i int) {
		actions = append(actions, fmt.Sprintf("foreach(%v)", i))
	})

	assert.Equal(t, []string{
		"filter(hello)",
		"map(hello)",
		"foreach(5)",
		"filter(my)",
		"filter(friend)",
		"map(friend)",
		"foreach(6)",
	}, actions)
}

func TestIter(t *testing.T) {
	res := map[int]int{}
	expectedI := 0
	for i, n := range Of(2, 3, 4, 5, 6).Iter() {
		require.Equal(t, expectedI, i)
		expectedI++
		res[i] = n
	}
	assert.Equal(t, map[int]int{0: 2, 1: 3, 2: 4, 3: 5, 4: 6}, res)
}

func TestIter_MultipleInstances(t *testing.T) {
	str := Of(1, 2, 3)
	var indexes, elems []int
	for i, n := range str.Iter() {
		indexes = append(indexes, i)
		elems = append(elems, n)
	}
	assert.Equal(t, []int{1, 2, 3}, elems)
	assert.Equal(t, []int{0, 1, 2}, indexes)
	// check that a new Iter instance returns a fresh iterator from the start
	var indexes2, elems2 []int
	for i, n := range str.Iter() {
		indexes2 = append(indexes2, i)
		elems2 = append(elems2, n)
	}
	assert.Equal(t, []int{1, 2, 3}, elems2)
	assert.Equal(t, []int{0, 1, 2}, indexes2)
}

func TestIterCombination(t *testing.T) {
	res := map[int]int{}
	expectedI := 0
	for i, n := range Of(0, 1, 2, 3, 4, 5, 6, 7, 8).
		Filter(func(i int) bool {
			return i%2 == 0
		}).Skip(2).Iter() {
		require.Equal(t, expectedI, i)
		expectedI++
		res[i] = n
	}
	assert.Equal(t, map[int]int{0: 4, 1: 6, 2: 8}, res)
}

func TestSeq(t *testing.T) {
	var res []int
	for n := range Of(2, 3, 4, 5, 6).Seq() {
		res = append(res, n)
	}
	assert.Equal(t, []int{2, 3, 4, 5, 6}, res)
}

func TestSeqCombination(t *testing.T) {
	var res []int
	for n := range Of(0, 1, 2, 3, 4, 5, 6, 7, 8).
		Filter(func(i int) bool {
			return i%2 == 0
		}).Skip(2).Seq() {
		res = append(res, n)
	}
	assert.Equal(t, []int{4, 6, 8}, res)
}

func TestSeq_MultipleInstances(t *testing.T) {
	str := Of(1, 2, 3)
	var elems []int
	for n := range str.Seq() {
		elems = append(elems, n)
	}
	assert.Equal(t, []int{1, 2, 3}, elems)
	// check that a new Iter instance returns a fresh iterator from the start
	var elems2 []int
	for n := range str.Seq() {
		elems2 = append(elems2, n)
	}
	assert.Equal(t, []int{1, 2, 3}, elems2)
}
