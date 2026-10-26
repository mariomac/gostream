package examples

import (
	"cmp"
	"fmt"
	"math/rand"

	"github.com/mariomac/gostream/stream"
)

func ExampleGenerate() {
	fmt.Println("picking up 5 random numbers from higher to lower")
	stream.Generate(rand.Uint32).
		Limit(5).
		Sorted(func(l, r uint32) int {
			// inverts the order provided by cmp.Compare
			return -cmp.Compare(r, l)
		}).
		ForEach(func(n uint32) {
			fmt.Println(n)
		})
}
