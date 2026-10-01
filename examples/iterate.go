package main

import (
	"fmt"

	"github.com/mariomac/gostream/stream"
)

func main_iterate() {
	numbers := stream.Iterate(1, double).Limit(6)

	words := numbers.Map(asWord).ToSlice()

	fmt.Println(words)
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
