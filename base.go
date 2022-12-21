package stream

type Iterator[T any] interface {
	Next() (T, bool)
}

type IteratorErr[T any] interface {
	Iterator[T]
	Err() error
}

type noOpIteratorErr[T any] struct {
	Iterator[T]
}

func (_ noOpIteratorErr[T]) Err() error {
	return nil
}

type baseStream struct {
	source          IteratorErr[any]
	mapper          func(any) (any, error)
	filterPredicate func(any) bool
	takePredicate   func(any) bool

	extras map[string]any
}

func (s baseStream) Err() error {
	if s.extras["error"] != nil {
		return s.extras["error"].(error)
	}

	return nil

}

func (s baseStream) nextMap() (any, bool) {
	val, ok := s.source.Next()
	if !ok {
		return val, false
	}

	res, err := s.mapper(val)
	if err != nil {
		s.extras["error"] = err
		return res, false
	}

	return res, true
}

func (s baseStream) nextFilter() (any, bool) {
	for {
		val, ok := s.source.Next()
		if !ok {
			break
		}

		if s.filterPredicate(val) {
			return val, true
		}
	}

	return nil, false
}

func (s baseStream) nextTakeWhile() (any, bool) {
	val, ok := s.source.Next()
	if ok && s.takePredicate(val) {
		return val, true
	}
	return val, false
}

func (s baseStream) Next() (any, bool) {
	if s.Err() != nil {
		return nil, false
	}

	if s.mapper != nil {
		return s.nextMap()
	}

	if s.filterPredicate != nil {
		return s.nextFilter()
	}

	if s.takePredicate != nil {
		return s.nextTakeWhile()
	}

	return s.source.Next()
}

func (s baseStream) MapE(mapper func(any) (any, error)) baseStream {
	return baseStream{
		mapper: mapper,
		source: s,
		extras: s.extras,
	}
}

func (s baseStream) Filter(predicate func(any) bool) baseStream {
	return baseStream{
		filterPredicate: predicate,
		source:          s,
		extras:          s.extras,
	}
}

func (s baseStream) TakeWhile(predicate func(any) bool) baseStream {
	return baseStream{
		takePredicate: predicate,
		source:        s,
		extras:        s.extras,
	}
}

func (s baseStream) Limit(n int) baseStream {
	c := counterPredicates{
		n: n,
	}

	return s.TakeWhile(c.lessThan)
}

func (s baseStream) Skip(n int) baseStream {
	c := counterPredicates{
		n: n,
	}

	return s.Filter(c.moreOrEqualThan)
}
