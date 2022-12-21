package stream

import (
	"strings"
	"testing"

	"github.com/faiz-alhadiid/go-stream/testutils"
)

func TestStream(t *testing.T) {
	// High level test
	t.Run(`Methods: Of, FromSlice, Stream(Skip,Map,Filter,Limit,Slice) `, func(t *testing.T) {
		s := Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

		s = s.Skip(1)

		s = s.Map(func(i int) int {
			return i * i
		})

		s = s.Filter(func(i int) bool {
			return i%2 == 1
		}).Limit(3)

		res := s.Slice()

		expected := []int{9, 25, 49}

		testutils.AssertEqual(t, expected, res)
	})

	t.Run("Methods: FromChan, Zero, Stream(TakeWhile, FlatMap, Peek, Reduce)", func(t *testing.T) {
		ch := make(chan string, 8)
		strs := []string{"Hello", "World", "Go", "Stream Test", "", "hi"}

		for _, str := range strs {
			ch <- str
		}

		s := Zero[string]()

		_, ok := s.Next()

		testutils.AssertEqual(t, false, ok, "Zero must not return anything")

		s = FromChan(ch)

		s = s.TakeWhile(func(s string) bool {
			return len(s) != 0
		})

		s = s.Peek(func(s string) {})

		s = s.FlatMap(func(s string) Iterator[string] {
			var arr []string
			for _, c := range s {
				arr = append(arr, string([]rune{c}))
			}

			return FromSlice(arr)
		})

		result := s.Reduce("", func(acc, each string) (res string) {
			if strings.ToLower(each) == each {
				return acc
			}

			return acc + each
		})

		testutils.AssertEqual(t, "HWGST", result)

	})
}
