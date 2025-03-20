package main

import "testing"

func TestConsumePrefix(t *testing.T) {
	tests := []struct {
		input string
		num   int
		want  int
	}{
		{"#..", 1, 1},
		{"#...", 1, 1},
		{"?...", 1, 1},
		{"??...", 1, 2},
		{"..??...", 1, 2},
		{"..??...###", 1, 2},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cache := make(ConsumptionCache, 0)
			if got := consumePrefixOfStringSatisfyingBrokenSet(tt.input, tt.num, &cache); len(got) != tt.want {
				t.Errorf("consumePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindNumberOfMatchingSetsRecursive(t *testing.T) {
	tests := []struct {
		input string
		nums  []int
		want  int
	}{
		{"#..", []int{1}, 1},
		{"#...", []int{1}, 1},
		{"?...", []int{1}, 1},
		{"??...", []int{1}, 2},
		{"..??...", []int{1}, 2},
		{"..??...#", []int{1, 1}, 2},
		{"..?...###", []int{1, 3}, 1},
		{"?.?...###", []int{1, 1, 3}, 1},
		{"..??...###", []int{1, 1, 3}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cache := make(MatchingSetsCache)
			cache2 := make(ConsumptionCache)
			if got := findNumberOfMatchingSetsRecursive(tt.input, tt.nums, &cache, &cache2); got != tt.want {
				t.Errorf("findNumberOfMatchingSetsRecursive() = %v, want %v", got, tt.want)
			}
		})
	}
}
