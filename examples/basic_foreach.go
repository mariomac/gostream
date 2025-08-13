package main

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func main_basic_foreach() {
	stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11).
		Filter(isPrime).
		ForEach(func(n int) {
			fmt.Printf("%d is a prime number\n", n)
		})
}
