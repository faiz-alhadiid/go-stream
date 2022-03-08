package sequence

// Generator represents generator function
type Generator func(yield chan<- any)

// Generate genereatse Sequence from provided generator
func Generate(generator Generator) ISequence {
	seq := make(chan any)
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
func FromChan[T any](ch chan T) ISequence {
	return Generate(func(yield chan<- any) {
		for val := range ch {
			yield <- val
		}
		close(yield)
	})
}

// FromSlice returns Sequence from slice
func FromSlice[T any](slice []T) ISequence {
	return Generate(func(yield chan<- any) {
		for _, val := range slice {
			yield <- val
		}
		close(yield)
	})
}
