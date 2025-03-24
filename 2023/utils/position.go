package utils

import (
	"cmp"
	"fmt"
	"log"
)

type Position struct {
	X int
	Y int
}

func (p *Position) String() string {
	return fmt.Sprintf("(x=%d, y=%d)", p.X, p.Y)
}

func (p *Position) East(maxXValue int) (bool, Position) {
	if p.X < maxXValue {
		return true, Position{X: p.X + 1, Y: p.Y}
	}
	return false, Position{}
}

func (p *Position) West() (bool, Position) {
	if p.X > 0 {
		return true, Position{X: p.X - 1, Y: p.Y}
	}
	return false, Position{}
}

func (p *Position) North() (bool, Position) {
	if p.Y > 0 {
		return true, Position{X: p.X, Y: p.Y - 1}
	}
	return false, Position{}
}

func (p *Position) South(maxYValue int) (bool, Position) {
	if p.Y < maxYValue {
		return true, Position{X: p.X, Y: p.Y + 1}
	}
	return false, Position{}
}

func (p *Position) FollowDirection(d Direction, maxXValue int, maxYValue int) (bool, Position) {
	switch d {
	case North:
		return p.North()
	case East:
		return p.East(maxXValue)
	case West:
		return p.West()
	case South:
		return p.South(maxYValue)
	}
	log.Fatal("invalid direction", d)
	return false, Position{}
}

func (p *Position) ManhattanDistance(other Position) int {
	return Abs(p.X-other.X) + Abs(p.Y-other.Y)
}

func ComparePositions(p1 Position, p2 Position) int {
	return cmp.Or(cmp.Compare(p1.X, p2.X), cmp.Compare(p1.Y, p2.Y))
}
