package stream

func MapE[T, V any](it Iterator[T], mapper func(T) (V, error)) Stream[V] {
	s := FromIterator(it)

	mp := mapWrapper[T, V]{mapper}

	newStream := baseStream(s).MapE(mp.mapE)
	return Stream[V](newStream)
}

func Map[T, V any](it Iterator[T], mapper func(T) V) Stream[V] {
	me := noErrMapWrapper[T, V]{mapper}
	return MapE(it, me.mapE)
}

func ForEach[T any](it Iterator[T], consume func(val T)) {
	for {
		val, ok := it.Next()
		if !ok {
			break
		}

		consume(val)
	}
}

func Reduce[T, V any](it Iterator[T], initial V, reduce func(acc V, each T) V) (result V) {
	result = initial

	ForEach(it, func(val T) {
		result = reduce(result, val)
	})

	return result
}
