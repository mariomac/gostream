# Go Streams API

Type safe Stream processing library inspired in the [Java Streams API](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html).

[![Go Reference](https://pkg.go.dev/badge/github.com/mariomac/gostream.svg)](https://pkg.go.dev/github.com/mariomac/gostream)
[![Go Report Card](https://goreportcard.com/badge/github.com/mariomac/gostream)](https://goreportcard.com/report/github.com/mariomac/gostream)


## Table of contents

* [Requirements](#requirements)
* [Usage examples](#usage-examples)
* [Limitations](#limitations)
* [Performance](#performance)
* [Completion status](#completion-status)
* [Extra credits](#extra-credits)

## Requirements

* Go 1.24 or higher

This library makes intensive usage of [Type Parameters (generics)](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md) so it is not compatible with any Go version lower than 1.18.

## Usage examples

For more details about the API, and until [pkg.go.dev](https://pkg.go.dev/github.com/mariomac/gostream), 
is able to parse documentation for functions and types using generics, you can have a quick look
to the generated [go doc text descriptions](./docs), or just let that the embedded document browser of your
IDE does the job.

### Example 1: basic creation, transformation and iteration

1. Creates a literal stream containing all the integers from 1 to 11.
2. From the Stream, selects all the integers that are prime
3. Iterates the stream. For each filtered int, prints a message.

```go
import (
  "fmt"
  "github.com/mariomac/gostream/stream"
)

func main() {
  numbers := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)

  for _, n := range numbers.Filter(isPrime).Iter {
    fmt.Printf("%d is a prime number\n", n)
  }
}

func isPrime(n int) bool {
  for i := 2; i <= n/2; i++ {
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

Alternatively, you can use the `ForEach` method to iterate the stream in a functional way:

```go
func main() {
  stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11).
    Filter(isPrime).
    ForEach(func(n int) {
      fmt.Printf("%d is a prime number\n", n)
    })
}
```
### Example 2: generation, map, limit and slice conversion

1. Creates an **infinite** stream of random integers (no problem, streams are evaluated lazily!)
2. Divides the random integer to get a number between 1 and 6
3. Limits the infinite stream to 5 elements.
4. Collects the stream items as a slice.

```go
rnd := rand.New(rand.NewSource(time.Now().UnixMilli()))
fmt.Println("let me throw 5 times a dice for you")

results := stream.Generate(rnd.Int).
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

### Example 3: Generation from an iterator, Map to a different type

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

### Example 4: deduplication of elements

Following example requires to compare the elements of the Stream, so the Stream needs to be
composed by `comparable` elements (this is, accepted by the the `==` and `!=` operators):

1. Instantiate a `Stream` of `comparable` items.
2. Pass it to the `Distinct` method, that will return a copy of the original Stream without
   duplicates
3. Operating as any other stream.

```go
words := stream.Distinct(
  stream.Of("hello", "hello", "!", "ho", "ho", "ho", "!"),
).ToSlice()

fmt.Printf("Deduplicated words: %v\n", words)
```

Output:

```
Deduplicated words: [hello ! ho]
```

### Example 5: sorting from higher to lower

1. Generate a stream of uint32 numbers.
2. Picking up 5 elements.
3. Sorting them by the inverse natural order (from higher to lower)
   - It's **important** to limit the number of elements, avoiding invoking
     `Sorted` over an infinite stream (otherwise it would panic).

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

### Example 6: Reduce and helper functions

1. Generate an infinite incremental Stream (1, 2, 3, 4...) using the `stream.Iterate`
   instantiator and the `item.Increment` helper function.
2. Limit the generated to 8 elements
3. Reduce all the elements multiplying them using the item.Multiply helper function

```go
fac8, _ := stream.Iterate(1, item.Increment[int]).
    Limit(8).
    Reduce(item.Multiply[int])
fmt.Println("The factorial of 8 is", fac8)
```

Output: 

```
The factorial of 8 is 40320
```

## Limitations

Due to the initial limitations of Go generics, the API has the following limitations.
We will work on overcome them as long as new features are added to the Go type parameters
specification.

* You can use `Map` and `FlatMap` as method as long as the output element has the same type of the input.
  If you need to map to a different type, you need to use `stream.Map` or `stream.FlatMap` as functions.
* There is no `Distinct` method. There is only a `stream.Distinct` function.
* There is no `ToMap` method. There is only a `stream.ToMap` function.

## Performance

You might want to check: [Performance comparison of Go functional stream libraries](https://macias.info/entry/202212020000_go_streams.md).

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

If you want a more performant, parallelizable alternative to create data processing pipelines (following
a programming model focused on Extract-Transform-Load, ETL), you
could give a try to my alternative project: [PIPES: Processing In Pipeline-Embedded Stages](https://github.com/mariomac/pipes).

## Completion status

* Stream instantiation functions
  - [X] Comparable
  - [X] Concat
  - [X] Empty
  - [X] Generate
  - [X] Iterate
  - [X] Of
  - [X] OfMap
  - [x] OfSlice
  - [X] OfChannel
* Stream transformers
  - [X] Distinct
  - [X] Filter
  - [X] FlatMap
  - [X] Limit
  - [X] Map
  - [X] Peek
  - [X] Skip
  - [X] Sorted
* Collectors/Terminals
  - [X] ToMap
  - [X] ToSlice
  - [X] AllMatch
  - [X] AnyMatch
  - [X] Count
  - [X] FindFirst
  - [X] ForEach
  - [X] Max
  - [X] Min
  - [X] NoneMatch
  - [X] Reduce
  - [X] Iter
* Auxiliary Functions
  - [X] Add (for numbers)
  - [X] Increment (for numbers)
  - [X] IsZero
  - [X] Multiply (for numbers)
  - [X] Neg (for numbers)
  - [X] Not (for bools)
* Future
  - [ ] Collectors for future standard generic data structures
    - E.g. [ ] Join (for strings)
  - [ ] Allow users implement their own Comparable or Ordered types
  - [ ] More operations inspired in the Kafka Streams API
  - [ ] Parallel streams 
    - [ ] FindAny


## Extra credits

The Stream processing and aggregation functions are heavily inspired in the
[Java Stream Specification](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html).

Stream code documentation also used 
[Stream Javadoc](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html) as an
essential reference and might contain citations from it.

