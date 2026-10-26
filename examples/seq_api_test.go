package examples

import (
	"fmt"
	"slices"

	"github.com/mariomac/gostream/stream"
)

func ExampleStream_Seq() {
	goSeq := slices.Values([]int{1, 2, 3, 4, 5, 6, 7})

	// OfSeq allows viewing a standard Go iter.Seq as a Stream
	odds := stream.OfSeq(goSeq).Filter(func(i int) bool {
		return i%2 == 1
	})

	// .Seq method allows converting a Stream to a Go iter.Seq
	// so it plays well with the Go standard iterators' API
	for n := range odds.Seq() {
		fmt.Print(n, " ")
	}
	fmt.Println()
	// Output: 1 3 5 7
}
