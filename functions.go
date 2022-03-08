package sequence

import (
	"reflect"
)

func nilSafeValueOf(val any, t reflect.Type) reflect.Value {
	if val != nil {
		return reflect.ValueOf(val)
	}
	return reflect.Zero(t)
}

// Function represent function with exactly 1 input and 1 output
type Function = func(any) any

func castFunction(fun any) (Function, bool) {
	if fun == nil {
		return nil, false
	}
	if f, ok := fun.(Function); ok {
		return f, true
	}
	rvFun := reflect.ValueOf(fun)
	rtFun := rvFun.Type()
	if rtFun.Kind() != reflect.Func {
		return nil, false
	}

	if rtFun.NumIn() != 1 || rtFun.NumOut() != 1 {
		return nil, false
	}

	f := func(x any) any {
		xv := nilSafeValueOf(x, rtFun.In(0))
		res := rvFun.Call(vSlice(xv))
		return res[0].Interface()
	}
	return f, true

}

// Predicate represents function that have exactly 1 input and 1 output with bool as a type
type Predicate = func(any) bool

func castPredicate(fun any) (Predicate, bool) {
	if fun == nil {
		return nil, false
	}
	if f, ok := fun.(Predicate); ok {
		return f, true
	}

	rvFun := reflect.ValueOf(fun)
	rtFun := rvFun.Type()
	if rtFun.Kind() != reflect.Func {
		return nil, false
	}

	if rtFun.NumIn() != 1 || rtFun.NumOut() != 1 {
		return nil, false
	}
	if rtFun.Out(0).Kind() != reflect.Bool {
		return nil, false
	}

	f := func(x any) bool {
		xv := nilSafeValueOf(x, rtFun.In(0))
		res := rvFun.Call(vSlice(xv))
		return res[0].Bool()
	}
	return f, true

}

// Consumer represents function that have 1 input and no output
type Consumer = func(any)

func castConsumer(fun any) (Consumer, bool) {
	if fun == nil {
		return nil, false
	}
	if f, ok := fun.(Consumer); ok {
		return f, true
	}

	rvFun := reflect.ValueOf(fun)
	rtFun := rvFun.Type()
	if rtFun.Kind() != reflect.Func {
		return nil, false
	}

	if rtFun.NumIn() != 1 {
		return nil, false
	}

	f := func(x any) {
		xv := nilSafeValueOf(x, rtFun.In(0))
		rvFun.Call(vSlice(xv))
	}
	return f, true
}

// BiFunction represent Function with 2 input parameter
type BiFunction = func(any, any) any

func castBiFunction(fun any) (BiFunction, bool) {
	if fun == nil {
		return nil, false
	}
	if f, ok := fun.(BiFunction); ok {
		return f, true
	}
	rvFun := reflect.ValueOf(fun)
	rtFun := rvFun.Type()
	if rtFun.Kind() != reflect.Func {
		return nil, false
	}

	if rtFun.NumIn() != 2 || rtFun.NumOut() != 1 {
		return nil, false
	}

	f := func(x any, y any) any {
		xv := nilSafeValueOf(x, rtFun.In(0))
		yv := nilSafeValueOf(y, rtFun.In(1))
		res := rvFun.Call(vSlice(xv, yv))
		return res[0].Interface()
	}
	return f, true
}
