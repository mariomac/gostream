package main

import (
	"fmt"

	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/stream"
)

func main_reduce() {
	fac8, _ := stream.Iterate(1, item.Increment[int]).
		Limit(8).
		Reduce(item.Multiply[int])
	fmt.Println("The factorial of 8 is", fac8)
}
