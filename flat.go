package stream

type flatStream[T any] struct {
	Iterator[Iterator[T]]
	temp Iterator[T]
}

func (f *flatStream[T]) Next() (T, bool) {
	var zero T

	for {
		if f.temp == nil {
			temp, ok := f.Iterator.Next()
			if !ok {
				return zero, false
			}

			f.temp = temp
		}

		val, ok := f.temp.Next()
		if ok {
			return val, true
		} else {
			f.temp = nil
		}
	}
}

func Flatten[T any](it Iterator[Iterator[T]]) Stream[T] {
	f := &flatStream[T]{
		Iterator: it,
	}

	return FromIterator[T](f)
}

func FlatMapE[T, V any](it Iterator[T], mapper func(T) (Iterator[V], error)) Stream[V] {
	stream := MapE(it, mapper)
	return Flatten[V](stream)
}

func FlatMap[T, V any](it Iterator[T], mapper func(T) Iterator[V]) Stream[V] {
	mp := noErrMapWrapper[T, Iterator[V]]{mapper}
	return FlatMapE(it, mp.mapE)
}
