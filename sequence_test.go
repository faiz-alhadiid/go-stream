package sequence

import (
	"reflect"
	"testing"
)

func TestSequence_Slice(t *testing.T) {
	type args struct {
		sIns []interface{}
	}
	tests := []struct {
		name      string
		s         *Sequence
		args      args
		wantSlice interface{}
		wantErr   bool
	}{
		{
			name: "Test no args",
			s: FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
				Filter(func(x int) bool { return x%2 == 0 }).
				Map(func(x int) int { return x * x }).
				Take(3).(*Sequence),
			wantSlice: []interface{}{4, 16, 36},
		},
		{
			name: "Test with args",
			s: FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
				Filter(func(x int) bool { return x%2 == 0 }).
				Map(func(x int) int { return x * x }).
				Take(3).(*Sequence),
			args: args{
				sIns: []interface{}{[]int{}},
			},
			wantSlice: []int{4, 16, 36},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSlice, err := tt.s.Slice(tt.args.sIns...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sequence.Slice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSlice, tt.wantSlice) {
				t.Errorf("Sequence.Slice() = %v, want %v", gotSlice, tt.wantSlice)
			}
		})
	}
}

func TestSequence_Reduce(t *testing.T) {
	type args struct {
		initial interface{}
		reducer interface{}
	}
	tests := []struct {
		name       string
		s          *Sequence
		args       args
		wantResult interface{}
		wantErr    bool
	}{
		{
			name: "Test success",
			s: FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
				Filter(func(x int) bool { return x%2 == 0 }).
				Map(func(x int) int { return x * x }).
				Skip(2).(*Sequence),
			args: args{
				initial: 999,
				reducer: func(acc, x int) int {
					if acc < x {
						return acc
					}
					return x
				},
			},
			wantResult: 64,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.s.Reduce(tt.args.initial, tt.args.reducer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sequence.Reduce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Sequence.Reduce() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
