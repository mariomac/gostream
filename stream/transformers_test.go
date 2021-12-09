package stream

import (
	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/order"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
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
		Distinct(Of(1, 1, 2, 3, 3, 3, 4, 5, 1, 2, 3, 4, 5)).
			ToSlice(),
	)
}

func TestSort(t *testing.T) {
	assert.Equal(t,
		[]int{1, 1, 2, 3, 5, 6, 7, 8, 8},
		Of(1, 7, 8, 3, 2, 1, 5, 8, 6).
			Sorted(order.Natural[int]).ToSlice())
}

func TestFlapMap(t *testing.T) {
	generateCharSequence := func(in string) Stream[byte] {
		return OfSlice([]byte(in))
	}
	generateNillableCharSequence := func(in string) Stream[byte] {
		if len(in) == 0 {
			return nil
		}
		return OfSlice([]byte(in))
	}

	assert.Empty(t, FlatMap(Of[string](), generateCharSequence).ToSlice())

	chars := FlatMap(
		Of[string]("", "hello", "my", "", "friends!", ""),
		generateCharSequence).ToSlice()
	assert.Equal(t, []byte("hellomyfriends!"), chars)

	chars = FlatMap(
		Of[string]("", "hello", "my", "", "friends!", ""),
		generateNillableCharSequence).ToSlice()
	assert.Equal(t, []byte("hellomyfriends!"), chars)
}

func TestFlapMap_Method(t *testing.T) {
	incrementalStream := func(length int) Stream[int] {
		return Iterate(1, item.Increment[int]).Limit(length)
	}

	items := FlatMap(Of[int](3, 2, 1, 0, 4), incrementalStream).ToSlice()
	assert.Equal(t, []int{1, 2, 3, 1, 2, 1, 1, 2, 3, 4}, items)
}
