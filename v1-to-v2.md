
* `Stream[T]` was an interface in v1. Now is a struct.
    - Consequently, instantiator methods which returned `Stream[T]` in v1, return `*Stream[T]` in v2.
* For simplification and coincidence with the
