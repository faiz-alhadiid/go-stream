package sequence

import (
	"reflect"
	"testing"
)

func Test_castFunction(t *testing.T) {
	type args struct {
		fun interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Function
		want1   bool
		argFunc interface{}
	}{
		{
			name: "Is Function",
			args: args{
				fun: func(a interface{}) interface{} { return a },
			},
			want:    func(a interface{}) interface{} { return a },
			want1:   true,
			argFunc: "a",
		},
		{
			name: "Is Nil",
			args: args{
				fun: nil,
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Is Not a func",
			args: args{
				fun: 1,
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Wrong input or output size",
			args: args{
				fun: func() {},
			},
			want:  nil,
			want1: false,
		},

		{
			name: "Valid func",
			args: args{
				fun: func(a string) interface{} { return a },
			},
			want:    func(a interface{}) interface{} { return a },
			want1:   true,
			argFunc: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := castFunction(tt.args.fun)
			if got1 != tt.want1 {
				t.Errorf("castFunction() got1 = %v, want %v", got1, tt.want1)
			}

			if tt.want != nil {
				res, resWant := got(tt.argFunc), tt.want(tt.argFunc)
				if !reflect.DeepEqual(res, resWant) {
					t.Errorf("Function() got = %v, want %v", res, resWant)
				}
			}
		})
	}
}

func Test_castPredicate(t *testing.T) {
	type args struct {
		fun interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Predicate
		want1   bool
		argFunc interface{}
	}{
		{
			name: "Is Predicate",
			args: args{
				fun: func(a interface{}) bool { return a == "a" },
			},
			want:    func(a interface{}) bool { return a == "a" },
			want1:   true,
			argFunc: "a",
		},
		{
			name: "Is Nil",
			args: args{
				fun: nil,
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Is Not a func",
			args: args{
				fun: 1,
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Wrong input or output size",
			args: args{
				fun: func() {},
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Output type is not bool",
			args: args{
				fun: func(_ interface{}) string { return "" },
			},
			want:  nil,
			want1: false,
		},

		{
			name: "Valid func",
			args: args{
				fun: func(a string) bool { return a == "a" },
			},
			want:    func(a interface{}) bool { return a == "a" },
			want1:   true,
			argFunc: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := castPredicate(tt.args.fun)
			if got1 != tt.want1 {
				t.Errorf("castPredicate() got1 = %v, want %v", got1, tt.want1)
			}

			if tt.want != nil {
				res, resWant := got(tt.argFunc), tt.want(tt.argFunc)
				if !reflect.DeepEqual(res, resWant) {
					t.Errorf("Predicate() got = %v, want %v", res, resWant)
				}
			}
		})
	}
}

func Test_castBiFunction(t *testing.T) {
	type args struct {
		fun interface{}
	}
	tests := []struct {
		name     string
		args     args
		want     BiFunction
		want1    bool
		argFunc1 interface{}
		argFunc2 interface{}
	}{
		{
			name: "Is BiFunction",
			args: args{
				fun: func(a interface{}, _ interface{}) interface{} { return a },
			},
			want:     func(a interface{}, _ interface{}) interface{} { return a },
			want1:    true,
			argFunc1: "a",
		},
		{
			name: "Is Nil",
			args: args{
				fun: nil,
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Is Not a func",
			args: args{
				fun: 1,
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Wrong input or output size",
			args: args{
				fun: func() {},
			},
			want:  nil,
			want1: false,
		},

		{
			name: "Valid func",
			args: args{
				fun: func(a string, _ string) interface{} { return a },
			},
			want:     func(a interface{}, _ interface{}) interface{} { return a },
			want1:    true,
			argFunc1: "a",
			argFunc2: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := castBiFunction(tt.args.fun)
			if got1 != tt.want1 {
				t.Errorf("castBiFunction() got1 = %v, want %v", got1, tt.want1)
			}

			if tt.want != nil {
				res, resWant := got(tt.argFunc1, tt.argFunc2), tt.want(tt.argFunc1, tt.argFunc2)
				if !reflect.DeepEqual(res, resWant) {
					t.Errorf("BiFunction() got = %v, want %v", res, resWant)
				}
			}
		})
	}
}

func Test_castConsumer(t *testing.T) {
	type args struct {
		fun interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  Consumer
		want1 bool
	}{
		{
			name: "Is Consumer",
			args: args{
				fun: func(a interface{}) {},
			},
			want:  func(a interface{}) {},
			want1: true,
		},
		{
			name: "Is Nil",
			args: args{
				fun: nil,
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Is Not a func",
			args: args{
				fun: 1,
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Wrong input or output size",
			args: args{
				fun: func() {},
			},
			want:  nil,
			want1: false,
		},
		{
			name: "Valid func",
			args: args{
				fun: func(a string) interface{} { return nil },
			},
			want:  func(a interface{}) {},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := castConsumer(tt.args.fun)
			if got1 != tt.want1 {
				t.Errorf("castConsumer() got1 = %v, want %v", got1, tt.want1)
			}
			if got != nil {
				got(nil)
			}
		})
	}
}
