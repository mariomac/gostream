package main

import "fmt"
import _ "strconv"

type Iterator[T any] interface {
	Next() (T, bool)
}

type SliceIterator[T any] struct {
	Items []T
}

func (si *SliceIterator[T]) Next() (T, bool) {
	if len(si.Items) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	n := si.Items[0]
	si.Items = si.Items[1:]
	return n, true
}

type Stream[T any] struct {
	Input Iterator[T]
}

func Map[IT any, OT any](is Stream[IT], fn func(IT) OT) Stream[OT] {
	// todo: lazy operation
	var out []OT
	for in, ok := is.Input.Next(); ok ; in, ok = is.Input.Next() {
		out = append(out, fn(in))
	}
	return Stream[OT]{Input: &SliceIterator[OT]{Items: out}}
}

func ForEach[IT any](is Stream[IT], fn func(IT)) {
	// todo: lazy operation
	for in, ok := is.Input.Next(); ok; in, ok = is.Input.Next() {
		fn(in)
	}
}

func main() {
	strs := Stream[string]{Input: &SliceIterator[string]{Items: []string{"f", "ba", "baz", "baee"}}}
	ints := Map[string, int](strs, func(it string) int {
		return len(it)
	})
	ForEach(ints, func(i int) {
		fmt.Println("len -> ", i)
	})

}
