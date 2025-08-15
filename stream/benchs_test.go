package stream

import (
	"fmt"
	"testing"
)

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

func BenchmarkForEach(b *testing.B) {
	for n := 0; n < b.N; n++ {
		count := 0
		sum := 0
		Generate(func() int {
			c := count
			count++
			return c
		}).Limit(iterations).ForEach(func(num int) {
			sum += num
		})
		if sum != 4950 {
			fmt.Println(sum)
			b.FailNow()
		}
	}
}

func BenchmarkIter(b *testing.B) {
	for n := 0; n < b.N; n++ {
		count := 0
		sum := 0
		for _, num := range Generate(func() int {
			c := count
			count++
			return c
		}).Limit(iterations).Iter() {
			sum += num
		}
		if sum != 4950 {
			fmt.Println(sum)
			b.FailNow()
		}
	}
}

func BenchmarkSeq(b *testing.B) {
	for n := 0; n < b.N; n++ {
		count := 0
		sum := 0
		for num := range Generate(func() int {
			c := count
			count++
			return c
		}).Limit(iterations).Seq() {
			sum += num
		}
		if sum != 4950 {
			fmt.Println(sum)
			b.FailNow()
		}
	}
}

func BenchmarkIterSlice(b *testing.B) {
	for n := 0; n < b.N; n++ {
		count := 0
		sum := 0
		for _, num := range Generate(func() int {
			c := count
			count++
			return c
		}).Limit(iterations).ToSlice() {
			sum += num
		}
		if sum != 4950 {
			fmt.Println(sum)
			b.FailNow()
		}
	}
}
