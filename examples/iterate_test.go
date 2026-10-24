package examples

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func ExampleIterate() {
	numbers := stream.Iterate(1, double).Limit(6)

	words := numbers.Map(asWord).ToSlice()

	fmt.Println(words)
	// Output: [one two four eight many many]
}

func double(n int) int {
	return 2 * n
}

func asWord(n int) string {
	if n < 10 {
		return []string{
			"zero", "one", "two", "three", "four", "five",
			"six", "seven", "eight", "nine",
		}[n]
	} else {
		return "many"
	}
}
