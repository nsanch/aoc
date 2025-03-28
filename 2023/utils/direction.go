package utils

import "log"

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

func (d Direction) Reverse() Direction {
	switch d {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	case West:
		return East
	}
	log.Fatal("Invalid direction", int(d))
	return 0
}

func (d Direction) String() string {
	switch d {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	}
	log.Fatal("Invalid direction", int(d))
	return ""
}

// returns DeltaX, DeltaY
func (d Direction) Delta() (int, int) {
	switch d {
	case North:
		return 0, -1
	case East:
		return 1, 0
	case South:
		return 0, 1
	case West:
		return -1, 0
	}
	log.Fatal("Invalid direction", d)
	return 0, 0
}
