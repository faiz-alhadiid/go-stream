package sequence

import "reflect"

//  Create slice from any variadic parameter
func vSlice[T any](val ...T) []T {
	return val
}

func fRecover() {
	recover()
}
func safeSend(c chan<- any, value any) {
	defer fRecover()
	c <- value
}

func safeClose(c chan any) {
	defer fRecover()
	close(c)
}

func zero(val any) any {
	tp := reflect.TypeOf(val)
	if tp == nil {
		return nil
	}
	return reflect.Zero(tp)
}
