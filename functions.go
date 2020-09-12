package sequence

import (
	"reflect"
)

func nilSafeValueOf(val interface{}, t reflect.Type) reflect.Value {
	if val != nil {
		return reflect.ValueOf(val)
	}
	return reflect.Zero(t)
}

// Function represent function with exactly 1 input and 1 output
type Function = func(interface{}) interface{}

func castFunction(fun interface{}) (Function, bool) {
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

	f := func(x interface{}) interface{} {
		xv := nilSafeValueOf(x, rtFun.In(0))
		res := rvFun.Call([]reflect.Value{xv})
		return res[0].Interface()
	}
	return f, true

}

// Predicate represents function that have exactly 1 input and 1 output with bool as a type
type Predicate = func(interface{}) bool

func castPredicate(fun interface{}) (Predicate, bool) {
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

	f := func(x interface{}) bool {
		xv := nilSafeValueOf(x, rtFun.In(0))
		res := rvFun.Call([]reflect.Value{xv})
		return res[0].Bool()
	}
	return f, true

}

// Consumer represents function that have 1 input and no output
type Consumer = func(interface{})

func castConsumer(fun interface{}) (Consumer, bool) {
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

	f := func(x interface{}) {
		xv := nilSafeValueOf(x, rtFun.In(0))
		rvFun.Call([]reflect.Value{xv})
	}
	return f, true
}

// BiFunction represent Function with 2 input parameter
type BiFunction = func(interface{}, interface{}) interface{}

func castBiFunction(fun interface{}) (BiFunction, bool) {
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

	f := func(x interface{}, y interface{}) interface{} {
		xv := nilSafeValueOf(x, rtFun.In(0))
		yv := nilSafeValueOf(y, rtFun.In(1))
		res := rvFun.Call([]reflect.Value{xv, yv})
		return res[0].Interface()
	}
	return f, true
}
