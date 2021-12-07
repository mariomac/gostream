package main

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func main_comparable() {
	words := stream.Distinct(
		stream.Of("hello", "hello", "!", "ho", "ho", "ho", "!"),
	).ToSlice()

	fmt.Printf("Deduplicated words: %v\n", words)
}
