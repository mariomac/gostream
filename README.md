# GoStream

Type safe Stream processing library inspired in the [Java Streams API](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html).

## Requirements

* Go 1.18.

This library makes intensive usage of [Type Parameters (generics)](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md) so it is not compatible with any version lower to Go 1.18.

Until Go 1.18 stable is officially released, you can download the development version of Go 1.18 using [Gotip](https://pkg.go.dev/golang.org/dl/gotip):

```
go install golang.org/dl/gotip@latest
gotip download
alias go=gotip
```

## Usage examples

1. Creates a literal stream containing all the integers from 1 to 11.
2. From the Stream, selects all the integers that are prime
3. For each filtered int, prints a message.

```go
import (
	"fmt"
	"github.com/mariomac/gostream/str"
)

func main() {
	stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11).
		Filter(isPrime).
		ForEach(func(n int) {
			fmt.Printf("%d is a prime number\n", n)
		})
}

func isPrime(n int) bool {
	for i := 2; i < n/2; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
```

Output: 
```
1 is a prime number
2 is a prime number
3 is a prime number
5 is a prime number
7 is a prime number
11 is a prime number
```

## ROADMAP

* Stream instantiation functions
  - [ ] Empty
  - [ ] Generate
  - [ ] Iterate
  - [X] Of
  - [ ] OfMap
  - [x] OfSlice
  - [ ] OfChannel
* Stream-to-Stream functions
  - [ ] Concat
  - [ ] Distinct
  - [X] Filter
  - [ ] FlatMap
  - [ ] Limit
  - [X] Map
  - [ ] Peek
  - [ ] Skip
  - [ ] Sorted
* Collectors/Terminals
  - [ ] ToMap
  - [X] ToSlice
  - [ ] AllMatch
  - [ ] AnyMatch
  - [ ] Count
  - [ ] FindAny
  - [ ] FindFirst
  - [X] ForEach
  - [ ] Max
  - [ ] Min
  - [ ] NoneMatch
  - [ ] Reduce
* Auxiliary Functions
  - [ ] Add (for numbers)
  - [ ] IsNil
  - [ ] IsZeroValue (superset of IsNil)
  - [ ] Join (for strings)
  - [ ] Mul (for numbers)
  - [ ] Neg (for numbers or bools)
* Other
  - [ ] Parallel streams 


## Extra credits

The Stream processing and aggregation functions are heavily inspired in the
[Java Stream Specification](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html).

Stream code documentation also used 
[Stream Javadoc](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html) as an
essential reference and might contain citations from it.

