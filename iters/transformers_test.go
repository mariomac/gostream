package iters

import (
	"fmt"
	"iter"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	in := slices.Values([]int{1, 2, 3, 4, 5})
	double := Map(in, func(i int) int {
		return i * 2
	})
	assert.Equal(t, []int{2, 4, 6, 8, 10}, slices.Collect(double))

	strings := Map(in, strconv.Itoa)
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, slices.Collect(strings))
}

func TestForEach(t *testing.T) {
	var copy []int
	in := slices.Values([]int{1, 2, 3, 4, 5})
	ForEach(in, func(i int) {
		copy = append(copy, i)
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, copy)
}

func TestFilter(t *testing.T) {
	in := slices.Values([]int{1, 2, 3, 4, 5})
	odds := Filter(in, func(n int) bool {
		return n%2 == 1
	})
	assert.Equal(t, []int{1, 3, 5}, slices.Collect(odds))
	empty := Filter(in, func(_ int) bool {
		return false
	})
	assert.Empty(t, slices.Collect(empty))
	all := Filter(in, func(_ int) bool {
		return true
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, slices.Collect(all))
}

func TestLimit(t *testing.T) {
	count := 0
	items := Limit(7, Generate(func() int {
		count++
		return count
	}))

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, slices.Collect(items))
}

func TestDistinct(t *testing.T) {
	in := slices.Values([]int{1, 1, 2, 3, 3, 3, 4, 5, 1, 2, 3, 4, 5})
	assert.Equal(t,
		[]int{1, 2, 3, 4, 5},
		slices.Collect(Distinct(in)),
	)
}

func TestFlapMap(t *testing.T) {
	generateCharSequence := func(in string) iter.Seq[byte] {
		return slices.Values([]byte(in))
	}
	generateNillableCharSequence := func(in string) iter.Seq[byte] {
		if len(in) == 0 {
			return Empty[byte]()
		}
		return generateCharSequence(in)
	}

	assert.Empty(t, slices.Collect(FlatMap(Empty[string](), generateCharSequence)))

	chars := FlatMap(
		slices.Values([]string{"", "hello", "my", "", "friends!", ""}),
		generateCharSequence)
	assert.Equal(t, []byte("hellomyfriends!"), slices.Collect(chars))

	chars = FlatMap(
		slices.Values([]string{"", "hello", "my", "", "friends!", ""}),
		generateNillableCharSequence)
	assert.Equal(t, []byte("hellomyfriends!"), slices.Collect(chars))
}

func TestPeek(t *testing.T) {
	var actions []string
	in := slices.Values([]int{1, 2, 3, 4, 5, 6})
	in = Filter(in, func(n int) bool { return n%2 == 1 })
	in = Peek(in, func(n int) {
		actions = append(actions, fmt.Sprint("processed ", n))
	})
	// until we don't aggregate the stream, Peek is not invoked due to lazyness
	assert.Empty(t, actions)

	require.Equal(t, []int{1, 3, 5}, slices.Collect(in))
	assert.Equal(t, []string{
		"processed 1", "processed 3", "processed 5",
	}, actions)
}

func TestSkip(t *testing.T) {
	assert.Empty(t, slices.Collect(Skip(3, Empty[int]())))
	assert.Empty(t, slices.Collect(Skip(3, slices.Values([]int{1, 2}))))
	assert.Empty(t, slices.Collect(Skip(3, slices.Values([]int{1, 2, 3}))))
	assert.Equal(t, []int{4, 5, 6}, slices.Collect(Skip(3, slices.Values([]int{1, 2, 3, 4, 5, 6}))))
}
