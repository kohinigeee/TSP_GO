package slicecol

import (
	"reflect"
	"testing"
)

func TestRemoveFast(t *testing.T) {
	tests := []struct {
		name     string
		argArray []int
		removIdx int
		want     []int
	}{
		{name: "remove 1", argArray: []int{1, 2, 3, 4, 5}, removIdx: 1, want: []int{1, 5, 3, 4}},
		{name: "one element ", argArray: []int{1}, removIdx: 0, want: []int{}},
	}

	for _, tt := range tests {
		got := RemoveFast(tt.argArray, tt.removIdx)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %v, want: %v", got, tt.want)
		}
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name     string
		argArray []int
		removIdx int
		want     []int
	}{
		{name: "remove 1", argArray: []int{1, 2, 3, 4, 5}, removIdx: 1, want: []int{1, 3, 4, 5}},
		{name: "one element ", argArray: []int{1}, removIdx: 0, want: []int{}},
	}

	for _, tt := range tests {
		got := Remove(tt.argArray, tt.removIdx)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %v, want: %v", got, tt.want)
		}
	}
}

func TestReverse(t *testing.T) {
	test := []struct {
		name     string
		argArray []int
		startIdx int
		endIdx   int
		want     []int
	}{
		{name: "reverse 1", argArray: []int{1, 2, 3, 4, 5}, startIdx: 1, endIdx: 3, want: []int{1, 4, 3, 2, 5}},
		{name: "reverse 2", argArray: []int{1, 2, 3, 4, 5}, startIdx: 0, endIdx: 4, want: []int{5, 4, 3, 2, 1}},
		{name: "reverse 3", argArray: []int{1, 2, 3, 4, 5}, startIdx: 2, endIdx: 3, want: []int{1, 2, 4, 3, 5}},
		{name: "same idxs", argArray: []int{1, 2, 3, 4, 5}, startIdx: 1, endIdx: 1, want: []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range test {
		got := Reverse(tt.argArray, tt.startIdx, tt.endIdx)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %v, want: %v", got, tt.want)
		}
	}
}
