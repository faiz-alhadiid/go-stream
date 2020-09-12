package sequence

var (
	opMap       = "map"
	opFilter    = "filter"
	opTakeWhile = "take_while"
)

type operation struct {
	Name      string
	Map       Function
	Filter    Predicate
	TakeWhile Predicate
}

type limitTakeWhile struct {
	counter int64
	limit   int64
}

func (l *limitTakeWhile) Check(_ interface{}) bool {
	t := l.counter < l.limit
	l.counter++
	return t
}

type skipPredicate struct {
	counter int64
	skip    int64
}

func (s *skipPredicate) Check(_ interface{}) bool {
	t := s.counter > s.skip
	s.counter++
	return t
}
