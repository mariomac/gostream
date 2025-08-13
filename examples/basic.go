package main

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

func main_basic() {
	numbers := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)

	for _, n := range numbers.Filter(isPrime).Iter {
		fmt.Printf("%d is a prime number\n", n)
	}
}
