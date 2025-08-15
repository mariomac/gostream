package stream

import (
	"cmp"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mariomac/gostream/item"
)

func TestLazyOperation(t *testing.T) {
	var actions []string
	in := Of("hello", "my", "friend")
	filtered := in.Filter(func(s string) bool {
		actions = append(actions, "filter("+s+")")
		return strings.Contains(s, "e")
	})
	mapped := Map(filtered, func(s string) int {
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

func TestInfiniteStreamAssertion(t *testing.T) {
	testCases := []func(s Stream[int]){
		func(s Stream[int]) {
			s.Sorted(cmp.Compare[int])
		},
		func(s Stream[int]) {
			s.ToSlice()
		},
		func(s Stream[int]) {
			ToMap(Map(s, func(i int) item.Pair[int, int] {
				return item.Pair[int, int]{Key: i, Val: i}
			}))
		},
		func(s Stream[int]) {
			s.Reduce(item.Add[int])
		},
		func(s Stream[int]) {
			s.AnyMatch(item.IsZero[int])
		},
		func(s Stream[int]) {
			s.AllMatch(item.IsZero[int])
		},
		func(s Stream[int]) {
			s.NoneMatch(item.IsZero[int])
		},
		func(s Stream[int]) {
			s.Count()
		},
	}
	for i, operation := range testCases {
		t.Run(fmt.Sprint("testcase", i), func(t *testing.T) {
			finished := make(chan struct{})
			go func() {
				defer func() {
					if r := recover(); r != nil {
						close(finished)
					}
				}()
				// Try to operate in an infinite stream must throw a panic
				operation(Generate(rand.Int).
					Map(rand.Intn).
					Filter(func(t int) bool {
						return true
					}))
			}()
			select {
			case <-finished:
			// ok
			case <-time.After(5 * time.Second):
				require.Fail(t, "timeout while expecting test to panic")
			}
		})
	}
}

func TestLimitInfiniteStreamAssertion(t *testing.T) {
	testCases := []func(s Stream[int]){
		func(s Stream[int]) {
			s.Sorted(cmp.Compare[int])
		},
		func(s Stream[int]) {
			s.ToSlice()
		},
		func(s Stream[int]) {
			ToMap(Map(s, func(i int) item.Pair[int, int] {
				return item.Pair[int, int]{Key: i, Val: i}
			}))
		},
		func(s Stream[int]) {
			s.Reduce(item.Add[int])
		},
		func(s Stream[int]) {
			s.AnyMatch(item.IsZero[int])
		},
		func(s Stream[int]) {
			s.AllMatch(item.IsZero[int])
		},
		func(s Stream[int]) {
			s.NoneMatch(item.IsZero[int])
		},
		func(s Stream[int]) {
			s.Count()
		},
	}
	for i, operation := range testCases {
		t.Run(fmt.Sprint("testcase", i), func(t *testing.T) {
			finished := make(chan struct{})
			go func() {
				// You must limit an infinite stream before trying to operate on it
				operation(Generate(rand.Int).
					Map(rand.Intn).
					Limit(10))
				close(finished)
			}()
			select {
			case <-finished:
			// ok
			case <-time.After(5 * time.Second):
				require.Fail(t, "timeout while expecting test to finish")
			}
		})
	}
}

func TestIter(t *testing.T) {
	res := map[int]int{}
	expectedI := 0
	for i, n := range Iter(Of(2, 3, 4, 5, 6)) {
		require.Equal(t, expectedI, i)
		expectedI++
		res[i] = n
	}
	assert.Equal(t, map[int]int{0: 2, 1: 3, 2: 4, 3: 5, 4: 6}, res)
}

func TestIter_MultipleInstances(t *testing.T) {
	str := Of(1, 2, 3)
	var indexes, elems []int
	for i, n := range Iter(str) {
		indexes = append(indexes, i)
		elems = append(elems, n)
	}
	assert.Equal(t, []int{1, 2, 3}, elems)
	assert.Equal(t, []int{0, 1, 2}, indexes)
	// check that a new Iter instance returns a fres iterator from the start
	var indexes2, elems2 []int
	for i, n := range Iter(str) {
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
	for n := range Seq(str) {
		elems = append(elems, n)
	}
	assert.Equal(t, []int{1, 2, 3}, elems)
	// check that a new Iter instance returns a fres iterator from the start
	var elems2 []int
	for n := range Seq(str) {
		elems2 = append(elems2, n)
	}
	assert.Equal(t, []int{1, 2, 3}, elems2)
}

func TestSeq2(t *testing.T) {
	input := OfMap(map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6"}).
		Filter(func(i item.Pair[int, string]) bool {
			return i.Key%2 == 0
		})
	output := map[int]string{}
	for k, v := range Seq2(input) {
		output[k] = v
	}
	assert.Equal(t, map[int]string{2: "2", 4: "4", 6: "6"}, output)
}

func TestSeq2_MultipleInstances(t *testing.T) {
	str := Of(
		item.Pair[int, int]{Key: 3, Val: 1},
		item.Pair[int, int]{Key: 2, Val: 2},
		item.Pair[int, int]{Key: 1, Val: 3})
	var keys, vals []int
	for k, v := range Seq2(str) {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	assert.Equal(t, []int{3, 2, 1}, keys)
	assert.Equal(t, []int{1, 2, 3}, vals)
	// check that a new Iter instance returns a fres iterator from the start
	var keys2, vals2 []int
	for k, v := range Seq2(str) {
		keys2 = append(keys2, k)
		vals2 = append(vals2, v)
	}
	assert.Equal(t, []int{3, 2, 1}, keys2)
	assert.Equal(t, []int{1, 2, 3}, vals2)
}
