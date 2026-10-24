package examples

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func ExampleStream_Distinct() {
	words := stream.Of("hello", "hello", "!", "ho", "ho", "ho", "!").Distinct().ToSlice()

	fmt.Printf("Deduplicated words: %v\n", words)
	// Output: Deduplicated words: [hello ! ho]
}
