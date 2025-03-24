package utils

import (
	"slices"
	"testing"
)

func Test_GetInteriorPoints(t *testing.T) {
	type args struct {
		name string
		area []string
		path []Position
		want []Position
	}
	tests := []args{{
		name: "1",
		area: []string{"###", "#.#", "###"},
		path: []Position{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}, {X: 1, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 1}},
		want: []Position{{X: 1, Y: 1}}},
		{
			name: "2",
			area: []string{"#####", "#...#", "#..##", "#..#.", "####."},
			path: []Position{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 1}, {X: 4, Y: 2}, {X: 3, Y: 2}, {X: 3, Y: 3}, {X: 3, Y: 4}, {X: 2, Y: 4}, {X: 1, Y: 4}, {X: 0, Y: 4}, {X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}},
			want: []Position{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 1, Y: 3}, {X: 2, Y: 3}}}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := make(Grid, len(tt.area))
			for y, row := range tt.area {
				grid[y] = []rune(row)
			}
			got := GetInteriorPoints(&grid, tt.path)
			slices.SortFunc(got, ComparePositions)
			slices.SortFunc(tt.want, ComparePositions)
			if !slices.Equal(got, tt.want) {
				t.Errorf("GetInteriorPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}
