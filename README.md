# gostream
Stream processing stuff for Go

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

