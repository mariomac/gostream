
* `Stream[T]` was an interface in v1. Now is a struct.
    - Consequently, instantiator methods which returned `Stream[T]` in v1, return `*Stream[T]` in v2.
* Removed all the transformation functions and replaced by methods. For example 
  `stream.Map(parent_stream, func)` is replaced by `parent_stream.Map(func)`
* 