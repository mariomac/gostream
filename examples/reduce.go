package main

import (
	"fmt"
	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/stream"
)

func main() {
	fac8, _ := stream.Generate(item.Incremental(1)).
		Limit(8).
		Reduce(item.Multiply[int])
	fmt.Println("The factorial of 8 is", fac8)
}
