package utils

import (
	"log"
	"maps"
	"slices"
	"strings"
)

type PositionRange struct {
	underlying IntegerRangeWithAxis
	startPos   Position
}

func NewPositionRangeFromValues(startPos Position, d Direction, length int) PositionRange {
	if length == 0 {
		log.Panic("length must be greater than 0")
	}

	var ret PositionRange
	switch d {
	case North:
		// for north and west we need to flip the start/end to ensure that we can always add length to start to get to the end point.
		startPos = Position{X: startPos.X, Y: startPos.Y - length + 1}
		return PositionRange{underlying: IntegerRangeWithAxis{start: startPos.Y, length: length, axis: "y"}, startPos: startPos}
	case West:
		startPos = Position{X: startPos.X - length + 1, Y: startPos.Y}
		return PositionRange{underlying: IntegerRangeWithAxis{start: startPos.X, length: length, axis: "x"}, startPos: startPos}
	case South:
		return PositionRange{underlying: IntegerRangeWithAxis{start: startPos.Y, length: length, axis: "y"}, startPos: startPos}
	case East:
		return PositionRange{underlying: IntegerRangeWithAxis{start: startPos.X, length: length, axis: "x"}, startPos: startPos}
	}
	return ret
}

func (r *PositionRange) Add(x Position) bool {
	if r.underlying.Axis() == "x" && x.Y == r.startPos.Y {
		return r.underlying.Add(x.X)
	} else if r.underlying.Axis() == "y" && x.X == r.startPos.X {
		return r.underlying.Add(x.Y)
	}
	return false
}

func (r *PositionRange) Absorb(r2 PositionRange) bool {
	if r.underlying.Axis() == r2.underlying.Axis() {
		if (r.underlying.Axis() == "x" && r.startPos.Y == r2.startPos.Y) ||
			(r.underlying.Axis() == "y" && r.startPos.X == r2.startPos.X) {
			return r.underlying.Absorb(r2.underlying)
		}
	} else {
		// if axes are different and r is a point that is on
		// the other's line we can flip its axis and absorb.
		if r.underlying.Length() == 1 && (r.underlying.Axis() == "x" && r.startPos.X == r2.startPos.X ||
			r.underlying.Axis() == "y" && r.startPos.Y == r2.startPos.Y) {
			if r2.underlying.Axis() == "x" {
				r.underlying = *NewIntegerRangeWithAxis(r2.startPos.X, r2.underlying.length, r2.underlying.Axis())
			} else {
				r.underlying = *NewIntegerRangeWithAxis(r2.startPos.Y, r2.underlying.length, r2.underlying.Axis())
			}
			r.startPos = r2.startPos
			return r.underlying.Absorb(r2.underlying)
		} else if r2.underlying.Length() == 1 && (r2.underlying.Axis() == "x" && r.startPos.X == r2.startPos.X ||
			r2.underlying.Axis() == "y" && r.startPos.Y == r2.startPos.Y) {
			// we already contain this point.
			return true
		}
	}
	return false
}

func (r *PositionRange) AreaOfIntersection(other *PositionRange) int {
	if r.underlying.Axis() == other.underlying.Axis() {
		if (r.underlying.Axis() == "x" && r.startPos.Y == other.startPos.Y) ||
			(r.underlying.Axis() == "y" && r.startPos.X == other.startPos.X) {
			return r.underlying.AreaOfIntersection(other.underlying)
		}
		return 0
	} else {
		if (r.underlying.Length() == 1 || other.underlying.Length() == 1) &&
			((r.underlying.Axis() == "x" && r.startPos.X == other.startPos.X) ||
				(r.underlying.Axis() == "y" && r.startPos.Y == other.startPos.Y)) {
			return 1
		} else if r.underlying.Axis() == "x" && other.underlying.Axis() == "y" {
			if r.underlying.Contains(other.startPos.X) && other.underlying.Contains(r.startPos.Y) {
				return 1
			}
		} else if r.underlying.Axis() == "y" && other.underlying.Axis() == "x" {
			if r.underlying.Contains(other.startPos.Y) && other.underlying.Contains(r.startPos.X) {
				return 1
			}
		}
	}
	return 0
}

func (r *PositionRange) EnumerateAllPointsSlow() []Position {
	ret := make([]Position, 0)
	if r.underlying.Axis() == "x" {
		for i := 0; i < r.underlying.Length(); i++ {
			ret = append(ret, Position{X: r.underlying.Start() + i, Y: r.startPos.Y})
		}
	} else {
		for i := 0; i < r.underlying.Length(); i++ {
			ret = append(ret, Position{X: r.startPos.X, Y: r.underlying.Start() + i})
		}
	}
	return ret
}

func (r *PositionRange) Height() int {
	if r.underlying.Axis() == "y" {
		return r.underlying.Length()
	}
	return 1
}

func (r *PositionRange) String() string {
	return r.underlying.String()
}

func (r *PositionRange) Width() int {
	if r.underlying.Axis() == "x" {
		return r.underlying.Length()
	}
	return 1
}

type PositionRanges struct {
	ranges []PositionRange
}

func (rs *PositionRanges) Add(r2 PositionRange) {
	for _, r := range rs.ranges {
		if r.Absorb(r2) {
			return
		}
	}
	rs.ranges = append(rs.ranges, r2)
}

func (rs *PositionRanges) AddAll(other *PositionRanges) {
	for _, subrange := range other.ranges {
		rs.Add(subrange)
	}
}

func (rs *PositionRanges) AreaOfIntersection(other *PositionRange) int {
	count := 0
	for _, r := range rs.ranges {
		count += r.AreaOfIntersection(other)
	}
	return count
}

/*
	func (rs *PositionRanges) CleanAndRemoveDuplication() {
		rationalized := new(PositionRanges)
		for _, r := range rs.ranges {
			rationalized.AddAll(r.FlipToX())
		}
		rs.ranges = rationalized.ranges
	}
*/
func (rs *PositionRanges) NumPoints() int {
	count := 0
	interCount := 0
	xCount := 0
	yCount := 0
	for idx, r := range rs.ranges {
		count += r.Width() * r.Height()
		if r.underlying.Axis() == "x" {
			xCount += r.Width()
		} else {
			yCount += r.Height()
		}
		//fmt.Println("range", r, "count: ", count)
		for j := idx + 1; j < len(rs.ranges); j++ {
			r2 := rs.ranges[j]
			inter := r.AreaOfIntersection(&r2)
			if inter > 0 {
				//fmt.Printf("%v and %v intersect by %d\n", r, r2, inter)
				interCount += inter
			}
		}
	}

	//x := rs.EnumerateAllPointsSlow()
	//slices.SortFunc(x, ComparePositions)
	//fmt.Println(x)

	//fmt.Printf("total count %d, xCount %d, yCount %d, total intersection %d total slow %d\n", count, xCount, yCount, interCount, len(rs.EnumerateAllPointsSlow()))
	return count - interCount
}

func (rs *PositionRanges) EnumerateAllPointsSlow() []Position {
	ret := make([]Position, 0)
	for _, r := range rs.ranges {
		ret = slices.Concat(ret, r.EnumerateAllPointsSlow())
	}
	return slices.Collect(maps.Keys(MakeSetFromSlice(ret)))
}

func (rs *PositionRanges) String() string {
	sb := strings.Builder{}
	for _, r := range rs.ranges {
		sb.WriteString(r.String())
		sb.WriteRune('\n')
	}
	return sb.String()
}
