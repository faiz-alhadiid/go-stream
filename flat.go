package sequence

type flatSequence struct {
	sequences  chan ISequence
	isInfinite bool
}

func (s *flatSequence) out() chan any {
	return nil
}

func (s *flatSequence) Map(fun any) ISequence {
	seq := FromChan(s.out())
	return seq.Map(fun)
}

func (s *flatSequence) Filter(predicate any) ISequence {
	seq := FromChan(s.out())
	return seq.Filter(predicate)
}

// IsInfinite show if Sequence is potentially infinite
func (s *flatSequence) IsInfinite() bool {
	return s.isInfinite
}

// TakeWhile filter element while predicate returns true
func (s *flatSequence) TakeWhile(predicate any) ISequence {
	seq := FromChan(s.out())
	return seq.TakeWhile(predicate)
}

// Take takes n element from sequence
func (s *flatSequence) Take(n int64) ISequence {
	seq := FromChan(s.out())
	return seq.Take(n)
}

// Skip skip n element from sequence
func (s *flatSequence) Skip(n int64) ISequence {
	seq := FromChan(s.out())
	return seq.Skip(n)
}

// Run runs sequence pipeline and output it to channel
func (s *flatSequence) Run() <-chan any {
	return s.out()
}
func (s *flatSequence) Reduce(initial any, reducer any) (result any, err error) {
	seq := FromChan(s.out())
	return seq.Reduce(initial, reducer)
}
func (s *flatSequence) ForEach(consumer any) (err error) {
	seq := FromChan(s.out())
	return seq.ForEach(consumer)
}
func (s *flatSequence) Err() error {
	return nil
}
