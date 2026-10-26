package examples

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func ExampleStream_Reduce() {
	fac8, _ := stream.Iterate(1, func(n int) int { return n + 1 }).
		Limit(8).
		Reduce(func(a, b int) int { return a * b })
	fmt.Println("The factorial of 8 is", fac8)
	// Output: The factorial of 8 is 40320
}
