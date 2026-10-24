package examples

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func isPrime(n int) bool {
	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func ExampleStream_Filter() {
	numbers := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)

	for n := range numbers.Filter(isPrime).Seq() {
		fmt.Printf("%d is a prime number\n", n)
	}
	// Output:
	// 1 is a prime number
	// 2 is a prime number
	// 3 is a prime number
	// 5 is a prime number
	// 7 is a prime number
	// 11 is a prime number
}
