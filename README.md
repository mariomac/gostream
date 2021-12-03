# gostream
Stream processing stuff for Go

## ROADMAP

* Stream instantiation functions
  - [ ] Generate
  - [ ] Iterate
  - [x] OfSlice
  - [ ] OfMap
  - [X] Of
* Stream-to-Stream functions
  - [X] Map
  - [ ] FlatMap
  - [ ] Filter
  - [ ] Distinct
  - [ ] Limit
  - [ ] Peek
  - [ ] Skip
  - [ ] Sorted
  - [ ] Concat
* Collectors
  - [ ] Count
  - [ ] FindAny
  - [ ] FindFirst
  - [ ] AllMatch
  - [ ] AnyMatch
  - [ ] Max
  - [ ] Min
  - [ ] NoneMatch
  - [ ] AsSlice
  - [ ] AsMap
  - [ ] Reduce
* Auxiliary Functions
  - [ ] Add (for numbers)
  - [ ] Neg (for numbers or bools)
  - [ ] Mul (for numbers)
  - [ ] Join (for strings)
  - [ ] IsNil
  - [ ] IsZeroValue (superset of IsNil)
* Other
  - [X] ForEach
  - [ ] Parallel streams 