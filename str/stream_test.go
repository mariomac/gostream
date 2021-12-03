package str

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLazyOperation(t *testing.T) {
	var actions []string
	in := Of("hello", "my", "friend")
	filtered := Filter(in, func(s string) bool {
		actions = append(actions, "filter("+s+")")
		return strings.Contains(s, "e")
	})
	mapped := Map(filtered, func(s string) int {
		actions = append(actions, "map("+s+")")
		return len(s)
	})
	ForEach(mapped, func(i int) {
		actions = append(actions, fmt.Sprintf("foreach(%v)", i))
	})

	assert.Equal(t, []string{
		"filter(hello)",
		"map(hello)",
		"foreach(5)",
		"filter(my)",
		"filter(friend)",
		"map(friend)",
		"foreach(6)",
	}, actions)

}
