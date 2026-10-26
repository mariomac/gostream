# Go Streams API changelog

## v1.0.0
* Major **breaking change** to adopt Go 1.27 generic methods
  * Removed all the transformation functions and replaced by methods. For example
    `stream.Map(parent_stream, func)` is replaced by `parent_stream.Map(func)`
  * `Stream[T]` was an interface. Now is a struct.
    - Consequently, instantiator methods which returned `Stream[T]`, now return `*Stream[T]`.
* Removed `item.Pair[K, V]` and created `Stream2[K, V]`
* For simplicity, no more infinite check. Java Streams do neither.
* Removed predicates in "item" package.
* Removed "order" package.
* New instantiation method: `Zip`

## v0.10.1
* Fix: `.Item` and `.Seq` methods are not used as a function reference anymore, but directly invoked.

## v0.10.0
* Added `stream.OfSeq` and `stream.OfSeq2` instantiation functions.
* Deprecated `order.Natural` and removed `order.Int` in favor of standard library's `cmp.Compare`
* Removed `order.SortSlice` in favor of `slices.Sort...` family of standard library functions
* `ToSlice` method will return a nil slice if the stream is empty. Before, it returned an empty slice.

## v0.9.0

* Added `Stream.Iter` method that allows directly using the stream within a for ... range loop.
* Added `Stream.Seq` method that allows to create an iter.Seq from a stream.
* Added `Seq2` function to iterate map-based streams as an iter.Sec2

## v0.8.1

* Adapted code to work in final version of Go 1.18, where the `constraints` package has been moved
  to  `golang.org/x/exp/constraints` (thanks to [mpetavy](https://github.com/mariomac/gostream/pull/2)
  for its contribution)

## v0.8.0

* Initial set of types, functions and methods (see [API documents](./docs) and
  _Completion status_ section in [README.md](./README.md))