package main

import (
	"fmt"
)
import _ "strconv"
import "github.com/mariomac/gostream/str"

func main() {
	str.Slice([]string{"f", "ba", "baz", "baee"}).
		Map(func(s string) string {
			return "*" + s
		}).ForEach(func(s string) {
		fmt.Println(s)
	})

	/*ints := str.Map[string, int](strs, func(it string) int {
		return len(it)
	})
	str.ForEach(ints, func(i int) {
		fmt.Println("len -> ", i)
	})*/

}
