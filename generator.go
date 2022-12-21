package stream

type funcIterator[T any] func() (T, bool)

func (f funcIterator[T]) Next() (T, bool) {
	return f()
}

func fromIteratorAny(it Iterator[any]) baseStream {
	source, ok := it.(IteratorErr[any])
	if !ok {
		source = noOpIteratorErr[any]{it}
	}

	return baseStream{
		source: source,
		extras: make(map[string]any),
	}

}

type iteratorErrAny[T any] struct {
	Iterator[T]
}

func (it iteratorErrAny[T]) Next() (any, bool) {
	return it.Iterator.Next()
}

func (it iteratorErrAny[T]) Err() error {
	conv, ok := it.Iterator.(IteratorErr[T])
	if ok {
		return conv.Err()
	}

	return nil
}

func FromIterator[T any](it Iterator[T]) Stream[T] {
	baseStream := fromIteratorAny(iteratorErrAny[T]{it})

	return Stream[T](baseStream)
}

func Generate[T any](f func() (T, bool)) Stream[T] {
	return FromIterator[T](funcIterator[T](f))
}

type sliceIterator[T any] struct {
	arr []T
}

func (s *sliceIterator[T]) Next() (T, bool) {
	var zero T
	if len(s.arr) == 0 {
		return zero, false
	}

	v := s.arr[0]

	s.arr = s.arr[1:]
	return v, true
}

func FromSlice[T any](arr []T) Stream[T] {
	it := &sliceIterator[T]{arr}
	return FromIterator[T](it)
}

func Of[T any](vals ...T) Stream[T] {
	return FromSlice(vals)
}

func zeroIteratorFunc[T any]() (T, bool) {
	var zero T
	return zero, false
}

func Zero[T any]() Stream[T] {
	zIt := zeroIteratorFunc[T]
	return Generate(zIt)
}

func FromChan[T any](ch <-chan T) Stream[T] {
	return Generate(func() (T, bool) {
		val, ok := <-ch
		return val, ok
	})
}
