package stream

type mapWrapper[T, V any] struct {
	f func(T) (V, error)
}

func (m mapWrapper[T, V]) mapE(a any) (any, error) {
	v, _ := a.(T)
	return m.f(v)
}

type noErrMapWrapper[T, V any] struct {
	f func(T) V
}

func (m noErrMapWrapper[T, V]) mapE(a T) (V, error) {
	return m.f(a), nil
}

type predicateWrapper[T any] struct {
	f func(T) bool
}

func (p predicateWrapper[T]) test(a any) bool {
	v, _ := a.(T)
	return p.f(v)
}

type counterPredicates struct {
	n       int
	current int
}

func (c *counterPredicates) lessThan(_ any) bool {
	r := c.current < c.n
	c.current++
	return r
}

func (c *counterPredicates) moreOrEqualThan(_ any) bool {
	r := c.current >= c.n
	c.current++
	return r
}
