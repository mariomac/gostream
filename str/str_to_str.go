package str

func ForEach[IT any](is Stream[IT], fn func(IT)) {
	it := is.iterator()
	for in, ok := it.next(); ok; in, ok = it.next() {
		fn(in)
	}
}

func Map[IT any, OT any](is Stream[IT], fn func(IT) OT) Stream[OT] {
	return (&connectedStream[IT, OT]{
		appliedFn: fn,
		Input:     is,
	}).attach()
}
