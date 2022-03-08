package sequence

import (
	"fmt"
)

// Make sure Sequence implements ISequence interface
var _ ISequence = &Sequence{}

// Sequence ...
type Sequence struct {
	seq        chan any
	isClosed   bool
	operations []operation
	parent     *Sequence
	isInfinite bool
	err        error
}

func (s *Sequence) IsInfinite() bool {
	return s.isInfinite
}

func (s *Sequence) close() {
	defer func() {
		recover()
		s.isClosed = true
	}()
	if s.parent != nil {
		s.parent.close()
	}
	if !s.isClosed {
		safeClose(s.seq)
	}
}

// Err ...
func (s *Sequence) Err() error { return s.err }

// Map ...
func (s *Sequence) Map(function any) ISequence {
	f, ok := castFunction(function)
	if !ok {
		panic("Failed to cast function")
	}
	op := operation{
		Name: opMap,
		Map:  f,
	}
	s.operations = append(s.operations, op)
	return s
}

// Filter ...
func (s *Sequence) Filter(predicate any) ISequence {
	p, ok := castPredicate(predicate)
	if !ok {
		panic("Failed to cast predicate")
	}
	op := operation{
		Name:   opFilter,
		Filter: p,
	}
	s.operations = append(s.operations, op)
	return s
}

// TakeWhile ...
func (s *Sequence) TakeWhile(predicate any) ISequence {
	p, ok := castPredicate(predicate)
	if !ok {
		panic("Failed to cast predicate")
	}
	op := operation{
		Name:      opTakeWhile,
		TakeWhile: p,
	}
	s.operations = append(s.operations, op)
	return s
}

// Take ...
func (s *Sequence) Take(n int64) ISequence {
	l := &limitTakeWhile{limit: n}
	s.isInfinite = false
	return s.TakeWhile(l.Check)
}

// Skip ...
func (s *Sequence) Skip(n int64) ISequence {
	skip := &skipPredicate{skip: n}
	s.isInfinite = false
	return s.Filter(skip.Check)
}

func (s *Sequence) FlatMap(fun func(interface{}) ISequence) ISequence {
	ch := make(chan ISequence)
	go func() {
		for seq := range s.Run() {
			defer fRecover()
			ch <- fun(seq)
		}
		close(ch)
	}()
	return &flatSequence{
		sequences: ch,
	}
}

func (s *Sequence) run(out chan any) {
	var i int
	var op operation

	defer func() {
		close(out)
		s.close()
		if r := recover(); r != nil {
			s.err = fmt.Errorf("Error at operation %d name %s: %v", i, s.operations[i].Name, r)
		}
	}()

seq:
	for value := range s.seq {
		for i, op = range s.operations {
			switch op.Name {

			case opMap:
				value = op.Map(value)

			case opFilter:
				if !op.Filter(value) {
					continue seq
				}
			case opTakeWhile:
				if !op.TakeWhile(value) {
					s.close()
					break seq
				}
			}
		}
		out <- value
	}

}

// Run runs prepared Sequence
func (s *Sequence) Run() <-chan any {
	out := make(chan any)
	go s.run(out)
	return out
}

// Reduce ...
func (s *Sequence) Reduce(initial any, reducer any) (result any, err error) {
	result = zero(initial)
	defer func() {
		r := recover()
		s.close()
		if r != nil {
			err = fmt.Errorf("Reduce error: %v", r)
		} else if s.err != nil {
			err = fmt.Errorf("Reduce error: %w", s.err)
		}
	}()

	if s.IsInfinite() {
		return result, PotentiallyInfiniteError
	}

	b, ok := castBiFunction(reducer)
	if !ok {
		panic("Failed to cast reducer")
	}
	result = initial
	for each := range s.Run() {
		result = b(result, each)
	}
	return
}

// ForEach ...
func (s *Sequence) ForEach(consumer any) (err error) {
	defer func() {
		r := recover()
		s.close()
		if r != nil {
			err = fmt.Errorf("ForEach error: %v", r)
		} else if s.err != nil {
			err = fmt.Errorf("ForEach error: %w", s.err)
		}
	}()
	c, ok := castConsumer(consumer)
	if !ok {
		panic("Failed to cast consumer")
	}
	for each := range s.Run() {
		c(each)
	}
	return
}

// CollectSlice collect Sequence result into slice
func CollectSlice[T any](s ISequence) (arr []T, err error) {
	if s.IsInfinite() {
		return nil, PotentiallyInfiniteError
	}

	s.ForEach(func(o any) {
		arr = append(arr, o.(T))
	})
	return arr, nil
}

func CollectChan[T any](s ISequence, cap ...int) (c chan T, err error) {
	if len(cap) == 0 {
		c = make(chan T)
	} else {
		c = make(chan T, cap[0])
	}

	go func() {
		s.ForEach(func(o any) {
			c <- o.(T)
		})

		close(c)
	}()

	return c, nil
}
