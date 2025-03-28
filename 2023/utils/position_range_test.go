package utils

import (
	"testing"
)

func TestPositionRangeAdd(t *testing.T) {
	tests := []struct {
		name          string
		initialRange  PositionRange
		posToAdd      Position
		expected      bool
		expectedStart int
		expectedLen   int
	}{
		{
			name:          "add to horizontal range - valid position",
			initialRange:  NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			posToAdd:      Position{X: 8, Y: 10},
			expected:      true,
			expectedStart: 5,
			expectedLen:   4,
		},
		{
			name:          "add to horizontal range - position with wrong Y",
			initialRange:  NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			posToAdd:      Position{X: 8, Y: 11},
			expected:      false,
			expectedStart: 5,
			expectedLen:   3,
		},
		{
			name:          "add to vertical range - valid position",
			initialRange:  NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 3),
			posToAdd:      Position{X: 5, Y: 13},
			expected:      true,
			expectedStart: 10,
			expectedLen:   4,
		},
		{
			name:          "add to vertical range - position with wrong X",
			initialRange:  NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 3),
			posToAdd:      Position{X: 6, Y: 13},
			expected:      false,
			expectedStart: 10,
			expectedLen:   3,
		},
		{
			name:          "add already included position - horizontal",
			initialRange:  NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			posToAdd:      Position{X: 6, Y: 10},
			expected:      true,
			expectedStart: 5,
			expectedLen:   3,
		},
		{
			name:          "add already included position - vertical",
			initialRange:  NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 3),
			posToAdd:      Position{X: 5, Y: 11},
			expected:      true,
			expectedStart: 10,
			expectedLen:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initialRange.Add(tt.posToAdd)
			if result != tt.expected {
				t.Errorf("Add() = %v, want %v", result, tt.expected)
			}

			var actualStart, actualLen int
			if tt.initialRange.underlying.Axis() == "x" {
				actualStart = tt.initialRange.underlying.Start()
				actualLen = tt.initialRange.underlying.Length()
			} else {
				actualStart = tt.initialRange.underlying.Start()
				actualLen = tt.initialRange.underlying.Length()
			}

			if actualStart != tt.expectedStart || actualLen != tt.expectedLen {
				t.Errorf("After Add(): range = [%d, %d], want [%d, %d]",
					actualStart, actualLen, tt.expectedStart, tt.expectedLen)
			}
		})
	}
}

func TestPositionRangeAbsorb(t *testing.T) {
	tests := []struct {
		name           string
		range1         PositionRange
		range2         PositionRange
		expected       bool
		expectedAxis   string
		expectedStart  int
		expectedLength int
	}{
		{
			name:           "absorb horizontal ranges - adjacent",
			range1:         NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			range2:         NewPositionRangeFromValues(Position{X: 8, Y: 10}, East, 2),
			expected:       true,
			expectedAxis:   "x",
			expectedStart:  5,
			expectedLength: 5,
		},
		{
			name:           "absorb horizontal ranges - different Y",
			range1:         NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			range2:         NewPositionRangeFromValues(Position{X: 8, Y: 11}, East, 2),
			expected:       false,
			expectedAxis:   "x",
			expectedStart:  5,
			expectedLength: 3,
		},
		{
			name:           "absorb vertical ranges - adjacent",
			range1:         NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 3),
			range2:         NewPositionRangeFromValues(Position{X: 5, Y: 13}, South, 2),
			expected:       true,
			expectedAxis:   "y",
			expectedStart:  10,
			expectedLength: 5,
		},
		{
			name:           "absorb vertical ranges - different X",
			range1:         NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 3),
			range2:         NewPositionRangeFromValues(Position{X: 6, Y: 13}, South, 2),
			expected:       false,
			expectedAxis:   "y",
			expectedStart:  10,
			expectedLength: 3,
		},
		{
			name:           "absorb point into horizontal range",
			range1:         NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 1),
			range2:         NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 3),
			expected:       true,
			expectedAxis:   "y",
			expectedStart:  10,
			expectedLength: 3,
		},
		{
			name:           "absorb point into vertical range",
			range1:         NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 1),
			range2:         NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			expected:       true,
			expectedAxis:   "x",
			expectedStart:  5,
			expectedLength: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.range1.Absorb(tt.range2)
			if result != tt.expected {
				t.Errorf("Absorb() = %v, want %v", result, tt.expected)
			}

			if result {
				if tt.range1.underlying.Axis() != tt.expectedAxis {
					t.Errorf("After Absorb(): axis = %s, want %s",
						tt.range1.underlying.Axis(), tt.expectedAxis)
				}

				if tt.range1.underlying.Start() != tt.expectedStart || tt.range1.underlying.Length() != tt.expectedLength {
					t.Errorf("After Absorb(): range = [%d, %d], want [%d, %d]",
						tt.range1.underlying.Start(), tt.range1.underlying.Length(), tt.expectedStart, tt.expectedLength)
				}
			}
		})
	}
}

func TestPositionRangeAreaOfIntersection(t *testing.T) {
	tests := []struct {
		name     string
		range1   PositionRange
		range2   PositionRange
		expected int
	}{
		{
			name:     "horizontal ranges - no intersection",
			range1:   NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			range2:   NewPositionRangeFromValues(Position{X: 9, Y: 10}, East, 2),
			expected: 0,
		},
		{
			name:     "horizontal ranges - overlap by 1",
			range1:   NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			range2:   NewPositionRangeFromValues(Position{X: 7, Y: 10}, East, 2),
			expected: 1,
		},
		{
			name:     "vertical ranges - no intersection",
			range1:   NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 3),
			range2:   NewPositionRangeFromValues(Position{X: 5, Y: 14}, South, 2),
			expected: 0,
		},
		{
			name:     "vertical ranges - overlap by 1",
			range1:   NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 3),
			range2:   NewPositionRangeFromValues(Position{X: 5, Y: 12}, South, 2),
			expected: 1,
		},
		{
			name:     "different axis ranges - no intersection",
			range1:   NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			range2:   NewPositionRangeFromValues(Position{X: 6, Y: 11}, South, 2),
			expected: 0,
		},
		{
			name:     "different axis ranges - intersection at point",
			range1:   NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			range2:   NewPositionRangeFromValues(Position{X: 6, Y: 10}, South, 2),
			expected: 1,
		},
		{
			name:     "horizontal range and point - intersection",
			range1:   NewPositionRangeFromValues(Position{X: 5, Y: 10}, East, 3),
			range2:   NewPositionRangeFromValues(Position{X: 6, Y: 10}, South, 1),
			expected: 1,
		},
		{
			name:     "vertical range and point - intersection",
			range1:   NewPositionRangeFromValues(Position{X: 5, Y: 10}, South, 3),
			range2:   NewPositionRangeFromValues(Position{X: 5, Y: 11}, East, 1),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.range1.AreaOfIntersection(&tt.range2)
			if result != tt.expected {
				t.Errorf("AreaOfIntersection() = %v, want %v", result, tt.expected)
			}
		})
	}
}
