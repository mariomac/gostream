package stream

import "testing"

const iterations = 100

func BenchmarkImperative(b *testing.B) {
	for n := 0; n < b.N; n++ {
		count := 0
		var result []int
		for i := 0; i < iterations; i++ {
			if count%3 == 0 {
				result = append(result, count*count)
			}
			count++
		}
		_ = result
	}
}

func BenchmarkFunctional(b *testing.B) {
	for n := 0; n < b.N; n++ {
		count := 0
		_ = Generate(func() int {
			c := count
			count++
			return c
		}).Filter(func(n int) bool {
			return n%3 == 0
		}).Map(func(n int) int {
			return n * n
		}).Limit(iterations).ToSlice()
	}
}
