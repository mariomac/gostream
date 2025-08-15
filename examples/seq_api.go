package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"

	"github.com/mariomac/gostream/order"
	"github.com/mariomac/gostream/stream"
)

func main_seq_api() {
	// get a slice of os.DirEntry elements
	osFiles, err := os.ReadDir(".")
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	// Put osFiles slice into a stream and filter it
	justFiles := stream.OfSlice(osFiles).
		Filter(func(entry os.DirEntry) bool {
			return !entry.IsDir()
		})

	// 1. Sorting the stream by file size in descending order
	//    (using Sorted method cmp.Compare and order.Inverse)
	// 2. Limit the stream size to the top 3 files by size (using Limit method)
	sizeTop3 := justFiles.Sorted(order.Inverse(func(a, b os.DirEntry) int {
		ai, _ := a.Info()
		bi, _ := b.Info()
		return cmp.Compare(ai.Size(), bi.Size())
	})).Limit(3)

	// .Seq method allows iterating the stream within a for..range
	fmt.Println("Top 3 files:")
	for v := range sizeTop3.Seq() {
		fmt.Println(v.Name())
	}

	// .Iter method allows iterating the stream within an indexed for..range
	fmt.Println("Top 3 files (indexed):")
	for k, v := range sizeTop3.Iter() {
		fmt.Printf("[%v] %v\n", k, v.Name())
	}

	// .Seq method also allows connecting the stream to other Go
	// standard library functions that expect an iter.Seq input
	// for example, slices.Collect
	asSlice := slices.Collect(sizeTop3.Seq())
	fmt.Println("Top 3 files (as slice): ", asSlice)
}
