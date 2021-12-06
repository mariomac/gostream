package main

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func main_comparable() {
	words := stream.Comparing(
		stream.Of("hello", "hello", "!", "ho", "ho", "ho", "!"),
	).Distinct().ToSlice()

	fmt.Printf("Deduplicated words: %v\n", words)
}
