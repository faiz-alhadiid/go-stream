package sequence

// ISequence is an adapter for Sequence Operation
type ISequence interface {
	// IsInfinite show if Sequence is potentially infinite
	IsInfinite() bool
	// Map maps value with function
	Map(fun any) ISequence
	// Filter filter element with predicate
	Filter(predicate any) ISequence
	// TakeWhile filter element while predicate returns true
	TakeWhile(predicate any) ISequence
	// Take takes n element from sequence
	Take(n int64) ISequence
	// Skip skip n element from sequence
	Skip(n int64) ISequence
	// Run runs sequence pipeline and output it to channel
	Run() <-chan any
	Reduce(initial any, reducer any) (result any, err error)
	ForEach(consumer any) (err error)
	Err() error
}
