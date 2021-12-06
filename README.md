# GoStream

Type safe Stream processing library inspired in the [Java Streams API](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html).

## Requirements

* Go 1.18.

This library makes intensive usage of [Type Parameters (generics)](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md) so it is not compatible with any Go version lower than 1.18.

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
  "github.com/mariomac/gostream/stream"
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

Next:

1. Creates an **infinite** stream of random integers (no problem, streams are evaluated lazily!)
2. Divides the random integer to get a number between 1 and 6
3. Limits the infinite stream to 5 elements.
4. Collects the stream items as a slice.

```go
rand.Seed(time.Now().UnixMilli())
fmt.Println("let me throw 5 times a dice for you")

results := stream.Generate(rand.Int).
    Map(func(n int) int {
        return n%6 + 1
    }).
    Limit(5).
    ToSlice()

fmt.Printf("results: %v\n", results)
```

Output:
```
let me throw 5 times a dice for you
results: [3 5 2 1 3]
```

Next: 

1. Generates an infinite stream composed by `1`, `double(1)`, `double(double(1))`, etc...
   and cut it to 6 elements.
2. Maps the numbers' stream to a strings' stream. Because, at the moment,
   [go does not allow type parameters in methods](https://github.com/golang/go/issues/49085),
   we need to invoke the `stream.Map` function instead of the `numbers.Map` method
   because the contained type of the output stream (`string`) is different than the type of
   the input stream (`int`).
3. Converts the words stream to a slice and prints it.


```go
func main() {
    numbers := stream.Iterate(1, double).Limit(6)
    words := stream.Map(numbers, asWord).ToSlice()
    fmt.Println(words)
}

func double(n int) int {
    return 2 * n
}

func asWord(n int) string {
    if n < 10 {
        return []string{"zero", "one", "two", "three", "four", "five",
            "six", "seven", "eight", "nine"}[n]
    } else {
        return "many"
    }
}
```

Output:
```
[one two four eight many many]
```

Next example requires to compare the elements of the Stream, so the Stream needs to be
converted to a `ComparableStream[T]`, where T must be `comparable` (this is, defining)
the `==` and `!=` operators:

1. Instantiate a `Stream` and make sure it is `comparable` by passing it to the
   `stream.Comparing` function.
2. Invoke the `Distinct` method, that will return a copy of the original Stream without
   duplicates (this method is only available in the `ComparableStream` interface).
3. Operating as any other stream.

```go
words := stream.Comparing(
  stream.Of("hello", "hello", "!", "ho", "ho", "ho", "!"),
).Distinct().ToSlice()

fmt.Printf("Deduplicated words: %v\n", words)
```

Output:

```
Deduplicated words: [hello ! ho]
```

Sorting a slice from higher to lower:

1. Generate a stream of uint32 numbers.
2. Picking up 5 elements.
3. Sorting them by the inverse natural order (from higher to lower)
   - It's **important** to limit the number of elements, avoiding invoking
     `Sorted` over an infinite stream (it would panic).

```go
fmt.Println("picking up 5 random numbers from higher to lower")
stream.Generate(rand.Uint32).
    Limit(5).
    Sorted(order.Inverse(order.Natural[uint32])).
    ForEach(func(n uint32) {
        fmt.Println(n)
    })
```

Output:

```
picking up 5 random numbers from higher to lower
4039455774
2854263694
2596996162
1879968118
1823804162
```


## Performance

Streams aren't the fastest option. They are aimed for complex workflows where you can
sacrifice few microseconds for the sake of code organization and readability. Also disclaimer:
functional streams don't have to always be the most readable option.

The following results show the difference in performance for an arbitrary set of operations
in an imperative form versus the functional form using streams (see
[stream/benchs_test.go file](stream/benchs_test.go)):

```
$ gotip test -bench=. -benchmem  ./...
goos: darwin
goarch: amd64
pkg: github.com/mariomac/gostream/stream
cpu: Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz
BenchmarkImperative-4            2098518               550.6 ns/op          1016 B/op          7 allocs/op
BenchmarkFunctional-4             293095              3653 ns/op            2440 B/op         23 allocs/op
```

## Completion status

* Stream instantiation functions
  - [X] Comparable
  - [ ] Empty
  - [X] Generate
  - [X] Iterate
  - [X] Of
  - [ ] OfMap
  - [x] OfSlice
  - [ ] OfChannel
* Stream transformers
  - [X] Concat
  - [X] Distinct
  - [X] Filter
  - [ ] FlatMap
  - [X] Limit
  - [X] Map
  - [ ] Peek
  - [ ] Skip
  - [X] Sorted
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

