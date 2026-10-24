package examples

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func ExampleDistinct() {
	words := stream.Of("hello", "hello", "!", "ho", "ho", "ho", "!")
	dedupe := stream.Distinct(words)

	fmt.Printf("Deduplicated words: %v\n", dedupe.ToSlice())
	// Output: Deduplicated words: [hello ! ho]
}
