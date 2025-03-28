package utils

import (
	"testing"
)

func Test_GetInteriorPoints(t *testing.T) {
	type args struct {
		name string
		area []string
		path []Position
		want int
	}
	tests := []args{{
		name: "1",
		area: []string{"###", "#.#", "###"},
		path: []Position{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
		want: 1},
		{
			name: "2",
			area: []string{"#####", "#...#", "#..##", "#..#.", "####."},
			path: []Position{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 1}, {X: 4, Y: 2}, {X: 3, Y: 2}, {X: 3, Y: 3}, {X: 3, Y: 4}, {X: 2, Y: 4}, {X: 1, Y: 4}, {X: 0, Y: 4}, {X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}},
			want: 7}}

	for _, tt := range tests[0:1] {
		t.Run(tt.name, func(t *testing.T) {
			got := GetInteriorPoints(tt.path)
			if got.NumPoints() != tt.want {
				t.Errorf("GetInteriorPoints() = %v, want %v", got.NumPoints(), tt.want)
			}
		})
	}
}
