package sequence

import "reflect"

// Generator represents generator function
type Generator func(chan<- interface{})

// Generate genereatse Sequence from provided generator
func Generate(generator Generator) *Sequence {
	seq := make(chan interface{})
	s := &Sequence{
		seq: seq,
	}
	go func() {
		defer func() {
			recover()
			s.close()
		}()
		generator(seq)
	}()
	return s
}

// FromChan return sequence from channel
func FromChan(ch interface{}) *Sequence {
	rv := reflect.ValueOf(ch)
	if rv.Kind() != reflect.Chan {
		panic("ch is not a channel")
	}

	return Generate(func(yield chan<- interface{}) {
		for {
			v, ok := rv.Recv()
			if !ok {
				break
			}
			safeSend(yield, v)
		}
		close(yield)
	})
}

// FromSlice returns Sequence from slice ...
func FromSlice(slice interface{}) *Sequence {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		panic("slice is not a Slice")
	}
	return Generate(func(yield chan<- interface{}) {
		for i := 0; i < rv.Len(); i++ {
			yield <- rv.Index(i).Interface()
		}
		close(yield)
	})
}
