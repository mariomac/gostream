package examples

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func ExampleStream_Distinct() {
	words := stream.Of("hello", "hello", "!", "ho", "ho", "ho", "!")
	dedupe := words.Distinct(stream.Cmp)

	fmt.Printf("Deduplicated words: %v\n", dedupe.ToSlice())
	// Output: Deduplicated words: [hello ! ho]
}
