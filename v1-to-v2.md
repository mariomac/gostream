
* `Stream[T]` was an interface in v1. Now is a function.
* Removed all the transformation functions and replaced by methods. For example 
  `stream.Map(parent_stream, func)` is replaced by `parent_stream.Map(func)`
* For simplicity, we don't check anymore about infinite streams. This behavior
  works in line with the source of inspiration (Java streams API)