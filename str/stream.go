package str

type Stream[T any] interface {
	iterator() iterator[T]
	// connects an abstract stream with its implementor. Must be called after each stream
	// implementation
	attach() Stream[T]
	Map(fn func(T) T) Stream[T]
	ForEach(fn func(T))
	AsSlice() []T
}

type iterator[T any] interface {
	next() (T, bool)
}

type abstractStream[T any] struct {
	implementor Stream[T]
}

func (bs *abstractStream[T]) Map(fn func(T) T) Stream[T] {
	return Map[T, T](bs.implementor, fn)
}

func (bs *abstractStream[T]) ForEach(fn func(T)) {
	ForEach[T](bs.implementor, fn)
}

type sliceStream[T any] struct {
	abstractStream[T]
	Items []T
}

func (si *sliceStream[T]) attach() Stream[T] {
	si.abstractStream.implementor = si
	return si
}

func (si *sliceStream[T]) iterator() iterator[T] {
	return &sliceIterator[T]{Items: si.Items}
}

type sliceIterator[T any] struct {
	Items []T
}

func (si *sliceIterator[T]) next() (T, bool) {
	if len(si.Items) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	n := si.Items[0]
	si.Items = si.Items[1:]
	return n, true
}

type connectedStream[IN, OUT any] struct {
	abstractStream[OUT]
	appliedFn func(IN) OUT
	Input     Stream[IN]
}

func (si *connectedStream[IN, OUT]) attach() Stream[OUT] {
	si.abstractStream.implementor = si
	return si
}

func (si *connectedStream[IN, OUT]) iterator() iterator[OUT] {
	return &connectedIterator[IN, OUT]{
		Input:     si.Input.iterator(),
		appliedFn: si.appliedFn,
	}
}

type connectedIterator[IN, OUT any] struct {
	appliedFn func(IN) OUT
	Input     iterator[IN]
}

func (c *connectedIterator[IN, OUT]) next() (OUT, bool) {
	n, ok := c.Input.next()
	if !ok {
		var zeroVal OUT
		return zeroVal, ok
	}
	return c.appliedFn(n), true
}
