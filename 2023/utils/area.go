package utils

import (
	"fmt"
	"iter"
	"log"
	"maps"
	"strings"
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

type SparseGrid struct {
	pathMap map[Position]rune
	height  int
	width   int
}

func (g *SparseGrid) String() string {
	if g.height > 10000 || g.width > 10000 {
		return "Grid too large to print"
	}
	sb := strings.Builder{}
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			sb.WriteRune(g.ItemAt(Position{X: x, Y: y}))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (g *SparseGrid) UnorderedPositions() iter.Seq[Position] {
	return maps.Keys(g.pathMap)
}

func (g *SparseGrid) GetNextAlongDirection(p Position, d Direction) *Position {
	deltaX, deltaY := d.Delta()
	var closest *Position
	for pos := range g.pathMap {
		if deltaX == Sign(pos.X-p.X) && deltaY == Sign(pos.Y-p.Y) {
			if closest == nil || p.ManhattanDistance(pos) < p.ManhattanDistance(*closest) {
				closest = &pos
			}
		}
	}
	//fmt.Printf("closest to point %v along direction %s is point %v\n", p, d.String(), closest)
	return closest
}

func (g *SparseGrid) ItemAt(p Position) rune {
	if v, ok := g.pathMap[p]; ok {
		return v
	}
	return '.'
}

func (g *SparseGrid) Set(p Position, v rune) {
	if p.X >= g.width {
		g.width = p.X + 1
	}
	if p.Y >= g.height {
		g.height = p.Y + 1
	}
	g.pathMap[p] = v
}

func MakeSparseGridFromPath(path []Position) *SparseGrid {
	pathMap := make(map[Position]rune)
	ret := &SparseGrid{pathMap, 0, 0}
	for _, p := range path {
		ret.Set(p, '#')
	}
	return ret
}

func getGroundPointsInDirection(grid *SparseGrid, pos Position, dir Direction) *PositionRange {
	hasNeighbor, neighbor := pos.FollowDirection(dir, grid.width-1, grid.height-1)
	if !hasNeighbor {
		return nil
	}

	var ret PositionRange
	nextNonGround := grid.GetNextAlongDirection(pos, dir)
	if nextNonGround == nil {
		switch dir {
		case North:
			ret = NewPositionRangeFromValues(Position{X: pos.X, Y: 0}, South, pos.Y)
		case South:
			ret = NewPositionRangeFromValues(Position{X: pos.X, Y: pos.Y + 1}, South, grid.height-pos.Y-1)
		case East:
			_, east := pos.East(grid.width - 1)
			ret = NewPositionRangeFromValues(east, dir, grid.width-pos.X-1)
		case West:
			ret = NewPositionRangeFromValues(Position{X: 0, Y: pos.Y}, East, pos.X)
		}
	} else {
		if nextNonGround.ManhattanDistance(pos) == 1 {
			// no points to return since the very next point is also non-ground.
			return nil
		}
		ret = NewPositionRangeFromValues(neighbor, dir, max(Abs(nextNonGround.X-pos.X), Abs(nextNonGround.Y-pos.Y))-1)
	}
	//fmt.Println("getGroundPointsInDirection", pos.String(), dir.String(), ret.String())
	return &ret
}

func getGroundPointsInDirections(grid *SparseGrid, pos Position, dir []Direction) *PositionRanges {
	ret := new(PositionRanges)
	for _, d := range dir {
		r := getGroundPointsInDirection(grid, pos, d)
		if r != nil {
			ret.Add(*r)
		}
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

func GetInteriorPoints(path []Position) *PositionRanges {
	grid := MakeSparseGridFromPath(path)
	currPointOnShape := classifyPointOnShapeGivenThreePoints(path[len(path)-1], path[0], path[1])
	var lastPointOnShape EdgeKind

	var debugGrid *SparseGrid //= MakeSparseGridFromPath(path)
	if debugGrid != nil {
		debugGrid.Set(path[0], []rune(currPointOnShape.String())[0])
	}

	side1Direction, side2Direction := currPointOnShape.GetDirectionsOnEachSide()
	pointsOnSide1 := getGroundPointsInDirections(grid, path[0], side1Direction)
	pointsOnSide2 := getGroundPointsInDirections(grid, path[0], side2Direction)

	isFlipped := false

	fmt.Printf("grid height=%d, width=%d\n", grid.height, grid.width)

	// start at 1 since we already handled the first point above.
	for i := 1; i < len(path); i++ {
		pos := path[i]
		lastPointOnShape = currPointOnShape
		currPointOnShape = classifyPointOnShapeGivenThreePoints(path[i-1], path[i], path[(i+1)%len(path)])
		if debugGrid != nil {
			debugGrid.Set(pos, []rune(currPointOnShape.String())[0])
		}
		if shouldFlipSideAtTransition(lastPointOnShape, currPointOnShape) {
			//fmt.Println("Flipping sides")
			isFlipped = !isFlipped
		}
		if isFlipped {
			side2Direction, side1Direction = currPointOnShape.GetDirectionsOnEachSide()
		} else {
			side1Direction, side2Direction = currPointOnShape.GetDirectionsOnEachSide()
		}
		pointsOnSide1.AddAll(getGroundPointsInDirections(grid, pos, side1Direction))
		pointsOnSide2.AddAll(getGroundPointsInDirections(grid, pos, side2Direction))
		//fmt.Printf("side1: %v, side2: %v\n", pointsOnSide1, pointsOnSide2)
	}
	if debugGrid != nil {
		fmt.Println(debugGrid.String())
	}

	// Check along the border. Does side1 overlap with it? If so, return side2.
	//pointsOnSide1.CleanAndRemoveDuplication()
	//pointsOnSide2.CleanAndRemoveDuplication()

	leftBorder := NewPositionRangeFromValues(Position{X: 0, Y: 0}, South, grid.height)
	rightBorder := NewPositionRangeFromValues(Position{X: grid.width - 1, Y: 0}, South, grid.height)
	topBorder := NewPositionRangeFromValues(Position{X: 0, Y: 0}, East, grid.width)
	bottomBorder := NewPositionRangeFromValues(Position{X: 0, Y: grid.height - 1}, East, grid.width)
	if pointsOnSide1.NumPoints() == 0 ||
		pointsOnSide1.AreaOfIntersection(&leftBorder) > 0 ||
		pointsOnSide1.AreaOfIntersection(&rightBorder) > 0 ||
		pointsOnSide1.AreaOfIntersection(&topBorder) > 0 ||
		pointsOnSide1.AreaOfIntersection(&bottomBorder) > 0 {
		//fmt.Println(pointsOnSide2.NumPoints())
		//fmt.Println(len(MakeSetFromSlice(pointsOnSide2.EnumerateAllPointsSlow())))
		//fmt.Println(pointsOnSide2.String())
		return pointsOnSide2
	}

	//fmt.Println(pointsOnSide1.String())
	//fmt.Println(pointsOnSide1.NumPoints())
	//fmt.Println(len(MakeSetFromSlice(pointsOnSide1.EnumerateAllPointsSlow())))
	return pointsOnSide1
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func debugPrintSides(debugGrid *SparseGrid, side1Map map[Position]bool, side2Map map[Position]bool, grid *SparseGrid) string {
	for p := range grid.UnorderedPositions() {
		if side1Map[p] && side2Map[p] {
			debugGrid.Set(p, '3')
		} else if side1Map[p] {
			debugGrid.Set(p, '1')
		} else if side2Map[p] {
			debugGrid.Set(p, '2')
		}
	}
	return debugGrid.String()
}

func ShoelaceArea(path []Position) int {
	// Shoelace formula
	// https://en.wikipedia.org/wiki/Shoelace_theorem
	// https://math.stackexchange.com/questions/1218/how-to-calculate-the-area-of-a-polygon-given-its-vertices
	// https://www.cuemath.com/geometry/shoelace-theorem/
	area := 0
	for i := range path {
		j := (i + 1) % len(path)
		area += path[i].X * path[j].Y
		area -= path[j].X * path[i].Y
	}
	return Abs(area) / 2
}
