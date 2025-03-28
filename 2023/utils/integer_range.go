package utils

import (
	"fmt"
)

type IntegerRangeWithAxis struct {
	start  int
	length int
	axis   string
}

func NewIntegerRangeWithAxis(start int, length int, axis string) *IntegerRangeWithAxis {
	ret := new(IntegerRangeWithAxis)
	ret.start = start
	ret.length = length
	ret.axis = axis
	return ret
}

func (xr *IntegerRangeWithAxis) Start() int {
	return xr.start
}

func (xr *IntegerRangeWithAxis) Last() int {
	return xr.start + xr.length - 1
}

func (xr *IntegerRangeWithAxis) Length() int {
	return xr.length
}

func (xr *IntegerRangeWithAxis) Axis() string {
	return xr.axis
}

func (xr *IntegerRangeWithAxis) Add(x int) bool {
	if xr.length == 0 {
		xr.start = x
		xr.length = 1
		return true
	} else if x == xr.start-1 {
		xr.start = x
		xr.length++
		return true
	} else if x == xr.start+xr.length {
		xr.length++
		return true
	} else if xr.start <= x && x <= xr.Last() {
		return true
	}
	return false
}

func (xr *IntegerRangeWithAxis) AreaOfIntersection(pr IntegerRangeWithAxis) int {
	// this is imperfect since the lines can intersect at exactly one point but not a problem
	// for this class to solve.
	if xr.axis != pr.axis {
		return 0
	}

	// if pr is right before xr, potentially overlapping with it
	if pr.Start() <= xr.Start() && pr.Last() >= xr.Start() && pr.Last() <= xr.Last() {
		amountOfOverlap := pr.Last() - xr.Start() + 1
		return amountOfOverlap
	}

	// if pr starts within xr and ends after xr
	if pr.Start() >= xr.Start() && pr.Start() <= xr.Last() && pr.Last() >= xr.Last() {
		amountOfOverlap := xr.Last() - pr.Start() + 1
		return amountOfOverlap
	}

	// if pr is right after xr
	if xr.Last() < pr.Start() {
		return 0
	}

	// if pr is inside xr
	if xr.Start() <= pr.Start() && pr.Last() <= xr.Last() {
		return pr.Length()
	}

	// if xr is inside pr
	if pr.Start() <= xr.Start() && xr.Last() <= pr.Last() {
		return xr.Length()
	}

	return 0
}

func (xr *IntegerRangeWithAxis) Absorb(pr IntegerRangeWithAxis) bool {
	// this is imperfect since the lines can intersect at exactly one point but not a problem
	// for this class to solve.
	if xr.axis != pr.axis {
		return false
	}

	// if pr is right before xr, potentially overlapping with it
	if pr.Start() <= xr.Start() && pr.Last() >= xr.Start() && pr.Last() <= xr.Last() {
		amountOfOverlap := pr.Last() - xr.Start() + 1
		xr.start = pr.start
		xr.length += pr.length - amountOfOverlap
		return true
	}

	// if pr is right after xr
	if xr.Last()+1 == pr.Start() {
		xr.length += pr.Length()
		return true
	}

	// if pr is inside xr
	if xr.Start() <= pr.Start() && pr.Last() <= xr.Last() {
		return true
	}

	// if xr is inside pr
	if pr.Start() <= xr.Start() && xr.Last() <= pr.Last() {
		xr.start = pr.start
		xr.length = pr.length
		return true
	}

	return false
}

func (xr *IntegerRangeWithAxis) String() string {
	return fmt.Sprintf("start: %v, length: %d", xr.start, xr.length)
}

func (xr *IntegerRangeWithAxis) Contains(x int) bool {
	return xr.start <= x && x <= xr.Last()
}
