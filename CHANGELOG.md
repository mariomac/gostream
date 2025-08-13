# Go Streams API changelog

## v0.9.0

* Added `Stream.Iter` method that allows directly using the stream within a for ... range loop.
* Added `Stream.Seq` method that allows to create an iter.Seq from a stream.

## v0.8.1

* Adapted code to work in final version of Go 1.18, where the `constraints` package has been moved
  to  `golang.org/x/exp/constraints` (thanks to [mpetavy](https://github.com/mariomac/gostream/pull/2)
  for its contribution)

## v0.8.0

* Initial set of types, functions and methods (see [API documents](./docs) and
  _Completion status_ section in [README.md](./README.md))