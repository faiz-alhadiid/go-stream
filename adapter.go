package sequence

// ISequence is an adapter for Sequence Operation
type ISequence interface {
	Map(fun interface{}) ISequence
	Filter(predicate interface{}) ISequence
	TakeWhile(predicate interface{}) ISequence
	Take(n int64) ISequence
	Skip(n int64) ISequence
	Run() <-chan interface{}
	Reduce(initial interface{}, reducer interface{}) (result interface{}, err error)
	ForEach(consumer interface{}) (err error)
	Slice(sIns ...interface{}) (slice interface{}, err error)
	Err() error
}
