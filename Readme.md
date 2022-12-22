# Go-Stream
## Lazy Operated Functional Stream API

Example:
```go
sum := stream.Of[int](1,2,3,4,5,6).
        Filter(func(i int) bool {
            return i % 2 == 0
        }).
        Map(func( i int) bool {
            return i * i
        }).
        Limit(2).
        Reduce(0, func(acc int, each int) int) {
            return acc + each
        })
fmt.Println(sum) // 20 (4 + 16)
```

Currently, Go doesn't support type parameters in method. so code like `s.Map(func(int) string{})` cannot work. Here's the solution:

```go
s := stream.Of[string]("123", "321", "444")

sStr := Map[string, int](s, func (val string) int {
    ret, _ := strconv.Atoi(val)
    return ret
})

arr := sStr.Slice() // []int{123, 321, 444}
```

### List Of Available Methods

#### Generator Function:
- `Generate[T any](f func() (T, bool)) Stream[T]`
- `FromIterator[T any](f Iterator[T]) Stream[T]`
- `FromSlice[T any](arr []T Stream[T]`
- `Of[T any](vals ...T) Stream[T]`
- `FromChan[T any](c <-chan T) Stream[T]`
- `Zero[T any]() Stream[T]`

#### Operator Function:
- `MapE[T, V any](s Iterator[T], mapper func(T) (V, error)) Stream[V]`
- `Map[T, V any](s Iterator[T], mapper func(T) V) Stream[V]`
- `FlatMapE[T, V any](s Iterator[T], mapper func(T) (Iterator[V], error)) Stream[V]`
- `FlatMap[T, V any](s Iterator[T], mapper func(T) Iterator[V]) Stream[V]`
- `Flatten[T any, E Iterator[T]](s Iterator[E]) Stream[T]`
- `Reduce[T, V any](s Iterator[T], init V, reducer func(acc V, each T) V) V`

#### `Stream[T any]` methods:
- `Stream[T].Next() (T, bool)`
- `Stream[T].Err() error`
- `Stream[T].Map(func(T)T) Stream[T]`
- `Stream[T].MapE(func(T) (T, error)) Stream[T]`
- `Stream[T].Filter(func(T) bool) Stream[T]`
- `Stream[T].TakeWhile(func(T) bool) Stream[T]`
- `Stream[T].Limit(int) Stream[T]`
- `Stream[T].Skip(int) Stream[T]`
- `Stream[T].FlatMap(func(T) Iterator[T]) Stream[T]`
- `Stream[T].FlatMapE(func(T) (Iterator[T], error)) Stream[T]`
- `Stream[T].Reduce(init T, f func(acc, each T) T) T`
- `Stream[T].ForEach(f func(T))`
- `Stream[T].Slice() T[]`

