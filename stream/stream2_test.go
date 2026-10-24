package stream

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLazyOperation2(t *testing.T) {
	var actions []string
	in := OfMap(map[string]int{"hello": 1, "my": 2, "friend": 3})
	filtered := in.Filter(func(k string, v int) bool {
		actions = append(actions, fmt.Sprint("filter(", k, ",", v, ")"))
		return strings.Contains(k, "e")
	})
	mapped := filtered.Map(func(k string, v int) (int, int) {
		actions = append(actions, fmt.Sprint("map(", k, ",", v, ")"))
		return len(k), v
	})
	mapped.ForEach(func(k, v int) {
		actions = append(actions, fmt.Sprintf("foreach(%v,%v)", k, v))
	})

	assert.Equal(t, []string{
		"filter(hello,1)",
		"map(hello,1)",
		"foreach(5,1)",
		"filter(my,2)",
		"filter(friend,3)",
		"map(friend,3)",
		"foreach(6,3)",
	}, actions)
}

func TestSeq2(t *testing.T) {
	input := OfMap(map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6"}).
		Filter(func(k int, v string) bool {
			return k%2 == 0
		})
	output := map[int]string{}
	for k, v := range input.Seq2() {
		output[k] = v
	}
	assert.Equal(t, map[int]string{2: "2", 4: "4", 6: "6"}, output)
}

func TestSeq2_MultipleInstances(t *testing.T) {
	str := OfMap(map[int]int{3: 1, 2: 2, 1: 3})
	var keys, vals []int
	for k, v := range str.Seq2() {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	assert.Equal(t, []int{3, 2, 1}, keys)
	assert.Equal(t, []int{1, 2, 3}, vals)
	// check that a new Iter instance returns a fresh iterator from the start
	var keys2, vals2 []int
	for k, v := range str.Seq2() {
		keys2 = append(keys2, k)
		vals2 = append(vals2, v)
	}
	assert.Equal(t, []int{3, 2, 1}, keys2)
	assert.Equal(t, []int{1, 2, 3}, vals2)
}
