package stream

type Stream[T any] baseStream

func (s Stream[T]) Next() (T, bool) {
	v, ok := baseStream(s).Next()
	t, _ := v.(T)

	return t, ok
}

func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
	pr := predicateWrapper[T]{predicate}
	newStream := baseStream(s).Filter(pr.test)
	return Stream[T](newStream)
}

func (s Stream[T]) MapE(mapper func(T) (T, error)) Stream[T] {
	mp := mapWrapper[T, T]{mapper}
	newStream := baseStream(s).MapE(mp.mapE)
	return Stream[T](newStream)
}

func (s Stream[T]) Map(mapper func(T) T) Stream[T] {
	me := noErrMapWrapper[T, T]{mapper}
	return s.MapE(me.mapE)
}

func (s Stream[T]) FlatMapE(mapper func(T) (Iterator[T], error)) Stream[T] {
	str := FlatMapE[T](s, func(t T) (Iterator[any], error) {
		it, err := mapper(t)
		if err != nil {
			return nil, err
		}

		return Map(it, func(a T) any {
			return a
		}), nil
	})

	return Stream[T](str)
}

func (s Stream[T]) FlatMap(mapper func(T) Iterator[T]) Stream[T] {
	me := noErrMapWrapper[T, Iterator[T]]{mapper}
	return s.FlatMapE(me.mapE)
}

func (s Stream[T]) Peek(consume func(T)) Stream[T] {
	return s.Filter(func(t T) bool {
		consume(t)
		return true
	})
}

func (s Stream[T]) TakeWhile(predicate func(T) bool) Stream[T] {
	pr := predicateWrapper[T]{predicate}
	newStream := baseStream(s).TakeWhile(pr.test)
	return Stream[T](newStream)
}

func (s Stream[T]) Limit(n int) Stream[T] {
	newStream := baseStream(s).Limit(n)
	return Stream[T](newStream)
}

func (s Stream[T]) Skip(n int) Stream[T] {
	newStream := baseStream(s).Skip(n)
	return Stream[T](newStream)
}

func (s Stream[T]) ForEach(f func(value T)) {
	for {
		val, ok := s.Next()
		if !ok {
			break
		}

		f(val)
	}
}

func (s Stream[T]) Reduce(initial T, f func(acc T, each T) (res T)) (result T) {
	result = initial
	s.ForEach(func(t T) {
		result = f(result, t)
	})

	return result
}

func (s Stream[T]) Slice() []T {
	var res []T
	s.ForEach(func(t T) {
		res = append(res, t)
	})

	return res
}
