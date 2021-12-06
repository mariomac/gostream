package stream

import (
	"github.com/mariomac/gostream/order"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	in := Of[int](1, 2, 3, 4, 5)
	double := in.Map(func(i int) int {
		return i * 2
	})
	assert.Equal(t, []int{2, 4, 6, 8, 10}, double.ToSlice())

	strings := Map(in, strconv.Itoa)
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, strings.ToSlice())
}

func TestForEach(t *testing.T) {
	var copy []int
	Of[int](1, 2, 3, 4, 5).ForEach(func(i int) {
		copy = append(copy, i)
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, copy)
}

func TestFilter(t *testing.T) {
	in := Of[int](1, 2, 3, 4, 5)
	odds := in.Filter(func(n int) bool {
		return n%2 == 1
	})
	assert.Equal(t, []int{1, 3, 5}, odds.ToSlice())
	empty := in.Filter(func(_ int) bool {
		return false
	})
	assert.Empty(t, empty.ToSlice())
	all := in.Filter(func(_ int) bool {
		return true
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, all.ToSlice())
}

func TestLimit(t *testing.T) {
	count := 0
	items := Generate(func() int {
		count++
		return count
	}).Limit(7).ToSlice()

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, items)
}

func TestDistinct(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3, 4, 5},
		Comparing(Of(1, 1, 2, 3, 3, 3, 4, 5, 1, 2, 3, 4, 5)).
			Distinct().ToSlice(),
	)
}

func TestSort(t *testing.T) {
	assert.Equal(t,
		[]int{1, 1, 2, 3, 5, 6, 7, 8, 8},
		Of(1, 7, 8, 3, 2, 1, 5, 8, 6).
			Sorted(order.Natural[int]).ToSlice())
}

func TestSort_InfiniteStream(t *testing.T) {
	finished := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				close(finished)
			}
		}()
		// Try to sort in an infinite stream must throw a panic
		Comparing(Generate(rand.Int)).
			Map(rand.Intn).
			Filter(func(t int) bool {
				return true
			}).Sorted(order.Natural[int])
	}()
	select {
	case <-finished:
	// ok
	case <-time.After(5 * time.Second):
		require.Fail(t, "timeout while expecting test to finish")
	}
}

func TestSort_LimitInfiniteStream(t *testing.T) {
	finished := make(chan struct{})
	go func() {
		// You must limit an infinite stream before trying to sort it
		Comparing(Generate(rand.Int)).
			Map(rand.Intn).
			Limit(10).
			Sorted(order.Natural[int])
		close(finished)
	}()
	select {
	case <-finished:
	// ok
	case <-time.After(5 * time.Second):
		require.Fail(t, "timeout while expecting test to finish")
	}
}
