package utils

import (
	"testing"
)

func TestAdd(t *testing.T) {
	// Test adding to empty range
	emptyRange := NewIntegerRangeWithAxis(0, 0, "x")
	if !emptyRange.Add(5) {
		t.Error("Add() should return true when adding to empty range")
	}
	if emptyRange.Start() != 5 || emptyRange.Length() != 1 {
		t.Errorf("Expected range to be [5, 1], got [%d, %d]", emptyRange.Start(), emptyRange.Length())
	}

	// Test adding element before start
	range1 := NewIntegerRangeWithAxis(5, 3, "x")
	if !range1.Add(4) {
		t.Error("Add() should return true when adding element before start")
	}
	if range1.Start() != 4 || range1.Length() != 4 {
		t.Errorf("Expected range to be [4, 4], got [%d, %d]", range1.Start(), range1.Length())
	}

	// Test adding element after end
	range2 := NewIntegerRangeWithAxis(5, 3, "x")
	if !range2.Add(8) {
		t.Error("Add() should return true when adding element after end")
	}
	if range2.Start() != 5 || range2.Length() != 4 {
		t.Errorf("Expected range to be [5, 4], got [%d, %d]", range2.Start(), range2.Length())
	}

	// Test adding element within range
	range3 := NewIntegerRangeWithAxis(5, 3, "x")
	if !range3.Add(6) {
		t.Error("Add() should return true when adding element within range")
	}
	if range3.Start() != 5 || range3.Length() != 3 {
		t.Errorf("Expected range to remain [5, 3], got [%d, %d]", range3.Start(), range3.Length())
	}

	// Test adding element outside range
	range4 := NewIntegerRangeWithAxis(5, 3, "x")
	if range4.Add(9) {
		t.Error("Add() should return false when adding element outside range")
	}
	if range4.Start() != 5 || range4.Length() != 3 {
		t.Errorf("Expected range to remain [5, 3], got [%d, %d]", range4.Start(), range4.Length())
	}
}

func TestAreaOfIntersection(t *testing.T) {
	tests := []struct {
		name     string
		range1   *IntegerRangeWithAxis
		range2   *IntegerRangeWithAxis
		expected int
	}{
		{
			name:     "different axes",
			range1:   NewIntegerRangeWithAxis(5, 3, "x"),
			range2:   NewIntegerRangeWithAxis(6, 3, "y"),
			expected: 0,
		},
		{
			name:     "range2 partially overlaps range1 at start",
			range1:   NewIntegerRangeWithAxis(5, 3, "x"),
			range2:   NewIntegerRangeWithAxis(3, 3, "x"),
			expected: 1,
		},
		{
			name:     "range2 is right before range1",
			range1:   NewIntegerRangeWithAxis(5, 3, "x"),
			range2:   NewIntegerRangeWithAxis(1, 3, "x"),
			expected: 0,
		},
		{
			name:     "range2 is inside range1",
			range1:   NewIntegerRangeWithAxis(5, 5, "x"),
			range2:   NewIntegerRangeWithAxis(6, 2, "x"),
			expected: 2,
		},
		{
			name:     "range1 is inside range2",
			range1:   NewIntegerRangeWithAxis(6, 2, "x"),
			range2:   NewIntegerRangeWithAxis(5, 5, "x"),
			expected: 2,
		},
		{
			name:     "complete overlap",
			range1:   NewIntegerRangeWithAxis(5, 3, "x"),
			range2:   NewIntegerRangeWithAxis(5, 3, "x"),
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.range1.AreaOfIntersection(*tt.range2)
			if result != tt.expected {
				t.Errorf("AreaOfIntersection() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAbsorb(t *testing.T) {
	tests := []struct {
		name           string
		range1         *IntegerRangeWithAxis
		range2         *IntegerRangeWithAxis
		expected       bool
		expectedStart  int
		expectedLength int
	}{
		{
			name:           "different axes",
			range1:         NewIntegerRangeWithAxis(5, 3, "x"),
			range2:         NewIntegerRangeWithAxis(6, 3, "y"),
			expected:       false,
			expectedStart:  5,
			expectedLength: 3,
		},
		{
			name:           "range2 partially overlaps range1 at start",
			range1:         NewIntegerRangeWithAxis(5, 3, "x"),
			range2:         NewIntegerRangeWithAxis(3, 3, "x"),
			expected:       true,
			expectedStart:  3,
			expectedLength: 5,
		},
		{
			name:           "range2 is right after range1",
			range1:         NewIntegerRangeWithAxis(5, 3, "x"),
			range2:         NewIntegerRangeWithAxis(8, 2, "x"),
			expected:       true,
			expectedStart:  5,
			expectedLength: 5,
		},
		{
			name:           "range2 is inside range1",
			range1:         NewIntegerRangeWithAxis(5, 5, "x"),
			range2:         NewIntegerRangeWithAxis(6, 2, "x"),
			expected:       true,
			expectedStart:  5,
			expectedLength: 5,
		},
		{
			name:           "range1 is inside range2",
			range1:         NewIntegerRangeWithAxis(6, 2, "x"),
			range2:         NewIntegerRangeWithAxis(5, 5, "x"),
			expected:       true,
			expectedStart:  5,
			expectedLength: 5,
		},
		{
			name:           "no overlap or adjacency",
			range1:         NewIntegerRangeWithAxis(5, 3, "x"),
			range2:         NewIntegerRangeWithAxis(10, 3, "x"),
			expected:       false,
			expectedStart:  5,
			expectedLength: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.range1.Absorb(*tt.range2)
			if result != tt.expected {
				t.Errorf("Absorb() = %v, want %v", result, tt.expected)
			}
			if tt.range1.Start() != tt.expectedStart || tt.range1.Length() != tt.expectedLength {
				t.Errorf("After Absorb(): range = [%d, %d], want [%d, %d]",
					tt.range1.Start(), tt.range1.Length(), tt.expectedStart, tt.expectedLength)
			}
		})
	}
}
