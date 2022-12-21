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
        }
fmt.Println(sum) // 20 (4 + 16)
```

Currently, Go doesn't support type parameters in method. so code like `s.Map(func(int) string{})` cannot work. Here's the solution:

```go
s := Stream.Of[string]("123", "321", "444")

sStr := Map[string, int](s, func (val string) int {
    ret, _ := strconv.Atoi(val)
    return ret
})

arr := sStr.Slice() // []int{123, 321, 444}
```



