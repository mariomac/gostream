package stream

import (
	"cmp"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap2(t *testing.T) {
	in := Of[int](1, 2, 3, 4, 5)
	double := in.Map(func(i int) int {
		return i * 2
	})
	assert.Equal(t, []int{2, 4, 6, 8, 10}, double.ToSlice())

	strings := in.Map(strconv.Itoa)
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, strings.ToSlice())
}

func TestForEach2(t *testing.T) {
	var copy []int
	Of[int](1, 2, 3, 4, 5).ForEach(func(i int) {
		copy = append(copy, i)
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, copy)
}

func TestFilter2(t *testing.T) {
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

func TestLimit2(t *testing.T) {
	count := 0
	items := Generate2(func() (int, int) {
		count++
		return count, count * 10
	}).Limit(7).ToMap()

	assert.Equal(t, map[int]int{1: 10, 2: 20, 3: 30, 4: 40, 5: 50, 6: 60, 7: 70}, items)
}

func TestLimit2Boundaries(t *testing.T) {
	m := map[int]int{1: 10, 2: 20, 3: 30}
	assert.Empty(t, OfMap(m).Limit(0).ToMap())
	assert.Equal(t, m, OfMap(m).Limit(15).ToMap())
}

func TestDistinct2(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3, 4, 5},
		Of(1, 1, 2, 3, 3, 3, 4, 5, 1, 2, 3, 4, 5).Distinct(Cmp).ToSlice(),
	)
}

func TestSort2(t *testing.T) {
	assert.Equal(t,
		[]int{1, 1, 2, 3, 5, 6, 7, 8, 8},
		Of(1, 7, 8, 3, 2, 1, 5, 8, 6).
			Sorted(cmp.Compare[int]).ToSlice())
}

func TestFlapMap2(t *testing.T) {
	add := func(i byte) byte { return i + 1 }
	generateCharSequence := func(in string, _ byte) *Stream2[byte, byte] {
		return Zip(
			OfSlice([]byte(in)),
			Iterate(1, add),
		)
	}
	generateNillableCharSequence := func(in string, _ byte) *Stream2[byte, byte] {
		if len(in) == 0 {
			return nil
		}
		return Zip(
			OfSlice([]byte(in)),
			Iterate(1, add),
		)
	}

	assert.Empty(t, Empty2[string, byte]().FlatMap(generateCharSequence).ToMap())

	chars := Zip(
		Of[string]("", "hello", "my", "", "friends!", ""),
		Iterate[byte](1, add),
	).FlatMap(generateCharSequence).Keys().ToSlice()
	assert.Equal(t, []byte("hellomyfriends!"), chars)

	chars = Zip(
		Of[string]("", "hello", "my", "", "friends!", ""),
		Iterate[byte](1, add),
	).FlatMap(generateNillableCharSequence).Keys().ToSlice()
	assert.Equal(t, []byte("hellomyfriends!"), chars)
}

func TestPeek2(t *testing.T) {
	var actions []string

	sl := Of(1, 2, 3, 4, 5, 6).
		Filter(func(n int) bool {
			return n%2 == 1
		}).
		Peek(func(n int) {
			actions = append(actions, fmt.Sprint("processed ", n))
		}).
		ToSlice()

	require.Equal(t, []int{1, 3, 5}, sl)
	assert.Equal(t, []string{
		"processed 1", "processed 3", "processed 5",
	}, actions)
}

func TestSkip2(t *testing.T) {
	assert.Empty(t, Empty[int]().Skip(3).ToSlice())
	assert.Empty(t, Of(1, 2).Skip(3).ToSlice())
	assert.Empty(t, Of(1, 2, 3).Skip(3).ToSlice())
	assert.Equal(t, []int{4, 5, 6}, Of(1, 2, 3, 4, 5, 6).Skip(3).ToSlice())
}
