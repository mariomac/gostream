# gostream
Stream processing stuff for Go

## ROADMAP

* Stream instantiation functions
  - [ ] Generate
  - [ ] Iterate
  - [X] Of
  - [ ] OfMap
  - [x] OfSlice
* Stream-to-Stream functions
  - [ ] Concat
  - [ ] Distinct
  - [ ] Filter
  - [ ] FlatMap
  - [ ] Limit
  - [X] Map
  - [ ] Peek
  - [ ] Skip
  - [ ] Sorted
* Collectors/Terminals
  - [ ] AsMap
  - [X] AsSlice
  - [ ] AllMatch
  - [ ] AnyMatch
  - [ ] Count
  - [ ] FindAny
  - [ ] FindFirst
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
  - [X] ForEach
  - [ ] Parallel streams 