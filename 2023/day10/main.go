package main

import (
	"fmt"
	"log"
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

		toVisit = toVisit[1:]
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

func main() {
	fmt.Println(part1("day10-input-easy2.txt"))
	fmt.Println(part1("day10-input-easy.txt"))
	fmt.Println(part1("day10-input.txt"))

}
