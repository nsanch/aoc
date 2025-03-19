package main

import (
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"

	"github.com/nsanch/aoc/aoc2023/utils"
)

type GridWithPipes struct {
	grid utils.Grid
}

func (grid GridWithPipes) FindStartingPosition() utils.Position {
	for y, row := range grid.grid {
		for x, r := range row {
			if r == 'S' {
				return utils.Position{X: x, Y: y}
			}
		}
	}
	log.Fatal("No starting position found")
	return utils.Position{}
}

func (grid GridWithPipes) PipeAt(pos utils.Position) PipeKind {
	return PipeKind(grid.grid.ItemAt(pos))
}

type Graph map[utils.Position][]utils.Position

func (graph Graph) String() string {
	var sb strings.Builder
	for from, tos := range graph {
		sb.WriteString(fmt.Sprintf("%v -> %v\n", from, tos))
	}
	return sb.String()
}

func (graph Graph) AddEdge(from utils.Position, to utils.Position) {
	graph[from] = append(graph[from], to)
}

func (graph Graph) FindCycleBFS(from utils.Position) []utils.Position {
	visited := make(map[utils.Position]bool)
	pathsToNode := make(map[utils.Position][]utils.Position)
	toVisit := make([]utils.Position, 0)
	toVisit = append(toVisit, from)
	for len(toVisit) > 0 {
		curr := toVisit[len(toVisit)-1]

		pathToCurr := pathsToNode[curr]

		//log.Printf("Visiting %v via path %v ", curr, pathToCurr)

		toVisit = toVisit[:len(toVisit)-1]
		if curr == from && len(pathToCurr) > 0 {
			return pathToCurr
		}

		for _, neighbor := range graph[curr] {
			if visited[neighbor] && neighbor != from {
				// Don't go backwards along the path
				continue
			}

			toVisit = append(toVisit, neighbor)
			pathsToNode[neighbor] = make([]utils.Position, len(pathsToNode[curr]))
			copy(pathsToNode[neighbor], pathsToNode[curr])
			pathsToNode[neighbor] = append(pathsToNode[neighbor], curr)
		}

		visited[curr] = true
	}
	log.Fatal("No path found", graph)
	return nil
}

type PipeKind rune

func (pk PipeKind) IsCorner() bool {
	return pk == 'L' || pk == '7' || pk == 'J' || pk == 'F'
}

func (pk PipeKind) IsPipe() bool {
	return pk != '.'
}

func (pk PipeKind) ConnectsNorth() bool {
	return pk == '|' || pk == 'L' || pk == 'J' || pk == 'S'
}

func (pk PipeKind) ConnectsSouth() bool {
	return pk == '|' || pk == '7' || pk == 'F' || pk == 'S'
}

func (pk PipeKind) ConnectsEast() bool {
	return pk == '-' || pk == 'L' || pk == 'F' || pk == 'S'
}

func (pk PipeKind) ConnectsWest() bool {
	return pk == '-' || pk == '7' || pk == 'J' || pk == 'S'
}

func MakeGraphFromGrid(grid GridWithPipes) Graph {
	g := make(Graph)
	for y, row := range grid.grid {
		for x, ch := range row {
			currPos := utils.Position{Y: y, X: x}
			hasEast, east := currPos.East(len(row) - 1)
			hasWest, west := currPos.West()
			hasNorth, north := currPos.North()
			hasSouth, south := currPos.South(len(grid.grid) - 1)
			currPipe := PipeKind(ch)
			if currPipe.ConnectsNorth() && hasNorth && PipeKind(grid.grid.ItemAt(north)).ConnectsSouth() {
				g.AddEdge(currPos, north)
			}
			if currPipe.ConnectsEast() && hasEast && PipeKind(grid.grid.ItemAt(east)).ConnectsWest() {
				g.AddEdge(currPos, east)
			}
			if currPipe.ConnectsSouth() && hasSouth && PipeKind(grid.grid.ItemAt(south)).ConnectsNorth() {
				g.AddEdge(currPos, south)
			}
			if currPipe.ConnectsWest() && hasWest && PipeKind(grid.grid.ItemAt(west)).ConnectsEast() {
				g.AddEdge(currPos, west)
			}
		}
	}
	return g
}

func part1(fname string) int {
	grid := GridWithPipes{utils.ReadGridFromFile(fname)}
	graph := MakeGraphFromGrid(grid)
	startingPos := grid.FindStartingPosition()
	path := graph.FindCycleBFS(startingPos)
	return len(path) / 2
}

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

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

func getSidesOfPipe(pipe PipeKind) ([]Direction, []Direction) {
	switch pipe {
	case '|':
		return []Direction{West}, []Direction{East}
	case 'L':
		return []Direction{West, South}, []Direction{}
	case 'J':
		return []Direction{}, []Direction{South, East}
	case '7':
		return []Direction{}, []Direction{North, East}
	case 'F':
		return []Direction{North, West}, []Direction{}
	case '-':
		return []Direction{South}, []Direction{North}
	}
	log.Fatal("Invalid pipe", pipe)
	return nil, nil
}

func shouldFlipSideAtTransition(first, second PipeKind) bool {
	switch first {
	case '|':
		return false
	case 'L':
		return second == 'J'
	case 'J':
		return second == 'L' || second == '-'
	case '7':
		return second == 'F'
	case 'F':
		return second == '7' || second == '-'
	case '-':
		return second == 'F' || second == 'J'
	}
	log.Fatal("Invalid pipe", first)
	return false
}

func getGroundPointsInDirection(grid GridWithPipes, pos utils.Position, dir Direction) []utils.Position {
	var ret []utils.Position
	deltaX, deltaY := dir.Delta()
	for {
		pos.Y += deltaY
		pos.X += deltaX
		if pos.Y >= 0 && pos.Y < len(grid.grid) && pos.X >= 0 && pos.X < len(grid.grid[0]) {
			if grid.grid.ItemAt(pos) == '.' {
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

func getGroundPointsInDirections(grid GridWithPipes, pos utils.Position, dir []Direction) []utils.Position {
	var ret []utils.Position
	for _, d := range dir {
		ret = slices.Concat(ret, getGroundPointsInDirection(grid, pos, d))
	}
	return ret
}

func getInsidePoints(grid GridWithPipes, path []utils.Position) []utils.Position {
	currPipe := grid.PipeAt(path[0])
	var lastPipe PipeKind

	side1Direction, side2Direction := getSidesOfPipe(currPipe)
	pointsOnSide1 := getGroundPointsInDirections(grid, path[0], side1Direction)
	pointsOnSide2 := getGroundPointsInDirections(grid, path[0], side2Direction)

	isFlipped := false

	// start at 1 since we already handled the first point above.
	for i := 1; i < len(path); i++ {
		pos := path[i]
		lastPipe = currPipe
		currPipe = grid.PipeAt(pos)
		//fmt.Printf("Last pipe %c at %v. Curr Pipe %c at %v\n", lastPipe, path[i-1], currPipe, pos)
		if shouldFlipSideAtTransition(lastPipe, currPipe) {
			//fmt.Println("Flipping sides")
			isFlipped = !isFlipped
		}
		if isFlipped {
			side2Direction, side1Direction = getSidesOfPipe(currPipe)
		} else {
			side1Direction, side2Direction = getSidesOfPipe(currPipe)
		}

		pointsOnSide1 = slices.Concat(pointsOnSide1, getGroundPointsInDirections(grid, pos, side1Direction))
		pointsOnSide2 = slices.Concat(pointsOnSide2, getGroundPointsInDirections(grid, pos, side2Direction))
	}

	side1Map := make(map[utils.Position]bool)
	side2Map := make(map[utils.Position]bool)
	for _, p := range pointsOnSide1 {
		side1Map[p] = true
	}
	for _, p := range pointsOnSide2 {
		side2Map[p] = true
	}

	var sb strings.Builder
	for y, row := range grid.grid {
		for x, ch := range row {
			p := utils.Position{Y: y, X: x}
			if side1Map[p] && side2Map[p] {
				sb.WriteRune('3')
			} else if side1Map[p] {
				sb.WriteRune('1')
			} else if side2Map[p] {
				sb.WriteRune('2')
			} else {
				sb.WriteRune(ch)
			}
		}
		sb.WriteRune('\n')
	}
	fmt.Println(sb.String())

	for _, p := range pointsOnSide1 {
		// If any of the points on side 1 are on the edge, then we know that the inside is on side 2
		if p.X == 0 || p.Y == 0 || p.X == len(grid.grid[0])-1 || p.Y == len(grid.grid)-1 {
			return slices.Collect(maps.Keys(side2Map))
		}
	}
	return slices.Collect(maps.Keys(side1Map))
}

func IdentifyStartingPosPipe(graph Graph, grid GridWithPipes, startingPos utils.Position) PipeKind {
	startingPosNeighbors := graph[startingPos]
	var replacementForStartingPos rune
	_, north := startingPos.North()
	_, south := startingPos.South(len(grid.grid) - 1)
	_, east := startingPos.East(len(grid.grid[0]) - 1)
	_, west := startingPos.West()
	neighbor1 := startingPosNeighbors[0]
	neighbor2 := startingPosNeighbors[1]
	switch {
	case neighbor1 == north && neighbor2 == south || neighbor1 == south && neighbor2 == north:
		replacementForStartingPos = '|'
	case neighbor1 == east && neighbor2 == west || neighbor1 == west && neighbor2 == east:
		replacementForStartingPos = '-'
	case neighbor1 == north && neighbor2 == east || neighbor1 == east && neighbor2 == north:
		replacementForStartingPos = 'J'
	case neighbor1 == north && neighbor2 == west || neighbor1 == west && neighbor2 == north:
		replacementForStartingPos = 'L'
	case neighbor1 == south && neighbor2 == east || neighbor1 == east && neighbor2 == south:
		replacementForStartingPos = 'F'
	case neighbor1 == south && neighbor2 == west || neighbor1 == west && neighbor2 == south:
		replacementForStartingPos = '7'
	}

	return PipeKind(replacementForStartingPos)
}

func part2(fname string) int {
	grid := GridWithPipes{utils.ReadGridFromFile(fname)}
	graph := MakeGraphFromGrid(grid)
	startingPos := grid.FindStartingPosition()
	grid.grid[startingPos.Y][startingPos.X] = rune(IdentifyStartingPosPipe(graph, grid, startingPos))

	path := graph.FindCycleBFS(startingPos)
	nodesInPath := make(map[utils.Position]bool)
	for _, node := range path {
		nodesInPath[node] = true
	}
	// replace all nodes that aren't in path with .
	for y, row := range grid.grid {
		for x := range row {
			if !nodesInPath[utils.Position{X: x, Y: y}] {
				grid.grid[y][x] = '.'
			}
		}
	}

	inside := getInsidePoints(grid, path)
	fmt.Println(inside)
	return len(inside)
}

func main() {
	fmt.Println(part1("day10-input-easy2.txt"))
	fmt.Println(part1("day10-input-easy.txt"))
	fmt.Println(part1("day10-input.txt"))

	fmt.Println(part2("day10-input-easy3.txt"))
	fmt.Println(part2("day10-input-easy4.txt"))
	fmt.Println(part2("day10-input.txt"))
}
