package stream

type flatStream[T any, Inner Iterator[T]] struct {
	Iterator[Inner]
	temp Iterator[T]
}

func (f *flatStream[T, Inner]) Next() (T, bool) {
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

func Flatten[T any, Inner Iterator[T]](it Iterator[Inner]) Stream[T] {
	f := &flatStream[T, Inner]{
		Iterator: it,
	}

	return FromIterator[T](f)
}

func FlatMapE[T, V any](it Iterator[T], mapper func(T) (Iterator[V], error)) Stream[V] {
	stream := MapE(it, mapper)
	return Flatten[V, Iterator[V]](stream)
}

func FlatMap[T, V any](it Iterator[T], mapper func(T) Iterator[V]) Stream[V] {
	mp := noErrMapWrapper[T, Iterator[V]]{mapper}
	return FlatMapE(it, mp.mapE)
}
