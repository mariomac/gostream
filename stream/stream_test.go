package stream

import (
	"fmt"
	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/order"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			s.Sorted(order.Natural[int])
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
			s.Sorted(order.Natural[int])
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
