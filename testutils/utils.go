package testutils

import (
	"reflect"
)

type testingT interface {
	Errorf(format string, args ...any)
}

func AssertEqual(t testingT, expected, actual any, messages ...any) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("actual(%v) != expected(%v) %v", actual, expected, messages)
	}
}
