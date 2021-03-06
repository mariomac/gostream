package main

import (
	"fmt"
	"reflect"
	"runtime"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func main() {
	funcs := []func(){
		main_basic,
		main_comparable,
		main_dice,
		main_iterate,
		main_reduce,
		main_sort,
	}

	for _, f := range funcs {
		name := GetFunctionName(f)
		fmt.Printf("--- %s\n", name)
		f()
		fmt.Println()
	}
}
