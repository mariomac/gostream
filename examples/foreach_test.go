package examples

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func ExampleForEach() {
	stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11).
		Filter(isPrime).
		ForEach(func(n int) {
			fmt.Printf("%d is a prime number\n", n)
		})
	// Output:
	// 1 is a prime number
	// 2 is a prime number
	// 3 is a prime number
	// 5 is a prime number
	// 7 is a prime number
	// 11 is a prime number
}
