package main

import (
	"fmt"
	"github.com/mariomac/gostream/order"
	"math/rand"

	"github.com/mariomac/gostream/stream"
)

func main() {
	fmt.Println("picking up 5 random numbers from higher to lower")
	stream.Generate(rand.Uint32).
		Limit(5).
		Sorted(order.Inverse(order.Natural[uint32])).
		ForEach(func(n uint32) {
			fmt.Println(n)
		})
}
