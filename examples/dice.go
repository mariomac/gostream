package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mariomac/gostream/stream"
)

func main_dice() {
	rnd := rand.New(rand.NewSource(time.Now().UnixMilli()))
	fmt.Println("let me throw 5 times a dice for you")

	results := stream.Generate(rnd.Int).
		Map(func(n int) int {
			return n%6 + 1
		}).
		Limit(5).ToSlice()

	fmt.Printf("results: %v\n", results)
}
