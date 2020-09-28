package sequence

import (
	"fmt"
	"reflect"
)

// Make sure Sequence implements ISequence interface
var _ ISequence = &Sequence{}

// Sequence ...
type Sequence struct {
	seq        chan interface{}
	isClosed   bool
	operations []operation
	parent     *Sequence
	err        error
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
func (s *Sequence) Map(function interface{}) ISequence {
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
func (s *Sequence) Filter(predicate interface{}) ISequence {
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
func (s *Sequence) TakeWhile(predicate interface{}) ISequence {
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
	return s.TakeWhile(l.Check)
}

// Skip ...
func (s *Sequence) Skip(n int64) ISequence {
	skip := &skipPredicate{skip: n}
	return s.Filter(skip.Check)
}

func (s *Sequence) run(out chan interface{}) {
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
func (s *Sequence) Run() <-chan interface{} {
	out := make(chan interface{})
	go s.run(out)
	return out
}

// Reduce ...
func (s *Sequence) Reduce(initial interface{}, reducer interface{}) (result interface{}, err error) {
	defer func() {
		r := recover()
		s.close()
		if r != nil {
			err = fmt.Errorf("Reduce error: %v", r)
		} else if s.err != nil {
			err = fmt.Errorf("Reduce error: %w", s.err)
		}
	}()
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
func (s *Sequence) ForEach(consumer interface{}) (err error) {
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

// Slice returns slice
func (s *Sequence) Slice(sIns ...interface{}) (slice interface{}, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("Slice error: %v", r)
		} else if s.err != nil {
			err = fmt.Errorf("Slice error: %w", s.err)
		}
	}()

	if len(sIns) == 0 {
		var arr []interface{}
		err = s.ForEach(func(val interface{}) {
			arr = append(arr, val)
		})
		slice = arr
		return
	}

	rvSlice := reflect.ValueOf(sIns[0])
	if rvSlice.Kind() != reflect.Slice {
		err = fmt.Errorf("Slice  error: param is not a slice")
		return
	}
	t := rvSlice.Type().Elem()
	err = s.ForEach(func(val interface{}) {
		rvSlice = reflect.Append(rvSlice, nilSafeValueOf(val, t))
	})

	if err != nil {
		return
	}

	slice = rvSlice.Interface()
	return
}
