package utils

import (
	"fmt"
	"log"
	"maps"
	"slices"
)

type EdgeKind int

const (
	ExteriorStraightNorthSouth EdgeKind = iota
	ExteriorStraightEastWest
	ExteriorCornerBottomLeft
	ExteriorCornerBottomRight
	ExteriorCornerTopRight
	ExteriorCornerTopLeft
)

// Stealing from PipeKind on day 10.
func (p EdgeKind) String() string {
	switch p {
	case ExteriorStraightNorthSouth:
		return "|"
	case ExteriorStraightEastWest:
		return "-"
	case ExteriorCornerBottomLeft:
		return "L"
	case ExteriorCornerBottomRight:
		return "J"
	case ExteriorCornerTopRight:
		return "7"
	case ExteriorCornerTopLeft:
		return "F"
	}
	log.Fatal("Invalid point on shape", int(p))
	return ""
}

func (p EdgeKind) GetDirectionsOnEachSide() ([]Direction, []Direction) {
	switch p {
	case ExteriorStraightNorthSouth:
		return []Direction{West}, []Direction{East}
	case ExteriorCornerBottomLeft:
		return []Direction{West, South}, []Direction{}
	case ExteriorCornerBottomRight:
		return []Direction{}, []Direction{South, East}
	case ExteriorCornerTopRight:
		return []Direction{}, []Direction{North, East}
	case ExteriorCornerTopLeft:
		return []Direction{North, West}, []Direction{}
	case ExteriorStraightEastWest:
		return []Direction{South}, []Direction{North}
	}
	log.Fatalf("Invalid point for GetDirectionsOnEachSide: %v", p.String())
	return nil, nil
}

func shouldFlipSideAtTransition(first, second EdgeKind) bool {
	switch first {
	case ExteriorStraightNorthSouth:
		return false
	case ExteriorCornerBottomLeft:
		return second == ExteriorCornerBottomRight
	case ExteriorCornerBottomRight:
		return second == ExteriorCornerBottomLeft || second == ExteriorStraightEastWest
	case ExteriorCornerTopRight:
		return second == ExteriorCornerTopLeft
	case ExteriorCornerTopLeft:
		return second == ExteriorCornerTopRight || second == ExteriorStraightEastWest
	case ExteriorStraightEastWest:
		return second == ExteriorCornerTopLeft || second == ExteriorCornerBottomRight
	}
	log.Fatal("Invalid point on shape for shouldFlipSideAtTransition", first.String())
	return false
}

func getGroundPointsInDirection(grid *Grid, pos Position, dir Direction) []Position {
	var ret []Position
	deltaX, deltaY := dir.Delta()
	for {
		pos.Y += deltaY
		pos.X += deltaX
		if pos.Y >= 0 && pos.Y < len(*grid) && pos.X >= 0 && pos.X < len((*grid)[0]) {
			if grid.ItemAt(pos) == '.' {
				ret = append(ret, pos)
			} else {
				break
			}
		} else {
			break
		}
	}
	return ret
}

func getGroundPointsInDirections(grid *Grid, pos Position, dir []Direction) []Position {
	var ret []Position
	for _, d := range dir {
		ret = slices.Concat(ret, getGroundPointsInDirection(grid, pos, d))
	}
	return ret
}

func classifyPointOnShapeGivenThreePoints(p0 Position, p1 Position, p2 Position) EdgeKind {
	type Delta struct {
		x int
		y int
	}
	zeroToOne := Delta{x: p1.X - p0.X, y: p1.Y - p0.Y}
	oneToTwo := Delta{x: p2.X - p1.X, y: p2.Y - p1.Y}
	switch zeroToOne {
	case Delta{x: 0, y: 1}:
		switch oneToTwo {
		case Delta{x: 1, y: 0}:
			return ExteriorCornerBottomLeft
		case Delta{x: -1, y: 0}:
			return ExteriorCornerBottomRight
		case Delta{x: 0, y: 1}:
			return ExteriorStraightNorthSouth
		}
		// case Delta{x: 0, y: -1}: is invalid because it's a reverse.
	case Delta{x: 0, y: -1}:
		switch oneToTwo {
		case Delta{x: 1, y: 0}:
			return ExteriorCornerTopLeft
		case Delta{x: -1, y: 0}:
			return ExteriorCornerTopRight
		case Delta{x: 0, y: -1}:
			return ExteriorStraightNorthSouth
		}
		// case Delta{x: 0, y: 1}: is invalid because it's a reverse.
	case Delta{x: 1, y: 0}:
		switch oneToTwo {
		case Delta{x: 0, y: 1}:
			return ExteriorCornerTopRight
		case Delta{x: 0, y: -1}:
			return ExteriorCornerBottomRight
		case Delta{x: 1, y: 0}:
			return ExteriorStraightEastWest
			// case Delta{x: -1, y: 0}: is invalid because it's a reverse.
		}
	case Delta{x: -1, y: 0}:
		switch oneToTwo {
		case Delta{x: 0, y: 1}:
			return ExteriorCornerTopLeft
		case Delta{x: 0, y: -1}:
			return ExteriorCornerBottomLeft
		case Delta{x: -1, y: 0}:
			return ExteriorStraightEastWest
			// case Delta{x: 1, y: 0}: is invalid because it's a reverse.
		}
	}
	log.Fatal("Invalid points for classifyPointOnShapeGivenThreePoints", p0.String(), p1.String(), p2.String())
	return ExteriorStraightNorthSouth
}

func GetInteriorPoints(grid *Grid, path []Position) []Position {
	currPointOnShape := classifyPointOnShapeGivenThreePoints(path[len(path)-1], path[0], path[1])
	var lastPointOnShape EdgeKind

	var debugGrid *Grid //grid.Clone()
	if debugGrid != nil {
		debugGrid.Set(path[0], []rune(currPointOnShape.String())[0])
	}

	side1Direction, side2Direction := currPointOnShape.GetDirectionsOnEachSide()
	pointsOnSide1 := getGroundPointsInDirections(grid, path[0], side1Direction)
	pointsOnSide2 := getGroundPointsInDirections(grid, path[0], side2Direction)

	isFlipped := false

	// start at 1 since we already handled the first point above.
	for i := 1; i < len(path); i++ {
		pos := path[i]
		lastPointOnShape = currPointOnShape
		currPointOnShape = classifyPointOnShapeGivenThreePoints(path[i-1], path[i], path[(i+1)%len(path)])
		if debugGrid != nil {
			debugGrid.Set(pos, []rune(currPointOnShape.String())[0])
		}
		//fmt.Printf("Last pipe %c at %v. Curr Pipe %c at %v\n", lastPipe, path[i-1], currPipe, pos)
		if shouldFlipSideAtTransition(lastPointOnShape, currPointOnShape) {
			//fmt.Println("Flipping sides")
			isFlipped = !isFlipped
		}
		if isFlipped {
			side2Direction, side1Direction = currPointOnShape.GetDirectionsOnEachSide()
		} else {
			side1Direction, side2Direction = currPointOnShape.GetDirectionsOnEachSide()
		}

		pointsOnSide1 = slices.Concat(pointsOnSide1, getGroundPointsInDirections(grid, pos, side1Direction))
		pointsOnSide2 = slices.Concat(pointsOnSide2, getGroundPointsInDirections(grid, pos, side2Direction))
	}
	if debugGrid != nil {
		fmt.Println(debugGrid.String())
	}

	side1Map := make(map[Position]bool)
	side2Map := make(map[Position]bool)
	for _, p := range pointsOnSide1 {
		side1Map[p] = true
	}
	for _, p := range pointsOnSide2 {
		side2Map[p] = true
	}

	if debugGrid != nil {
		fmt.Println(debugPrintSides(debugGrid, side1Map, side2Map, grid))
	}

	for _, p := range pointsOnSide1 {
		// If any of the points on side 1 are on the edge, then we know that the inside is on side 2
		if p.X == 0 || p.Y == 0 || p.X == len((*grid)[0])-1 || p.Y == len(*grid)-1 {
			return slices.Collect(maps.Keys(side2Map))
		}
	}
	return slices.Collect(maps.Keys(side1Map))
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func debugPrintSides(debugGrid *Grid, side1Map map[Position]bool, side2Map map[Position]bool, grid *Grid) string {
	for y, row := range *grid {
		for x := range row {
			p := Position{Y: y, X: x}
			if side1Map[p] && side2Map[p] {
				debugGrid.Set(p, '3')
			} else if side1Map[p] {
				debugGrid.Set(p, '1')
			} else if side2Map[p] {
				debugGrid.Set(p, '2')
			}
		}
	}
	return debugGrid.String()
}
