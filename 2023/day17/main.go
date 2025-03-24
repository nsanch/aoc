package main

import (
	"fmt"

	"github.com/nsanch/aoc/aoc2023/utils"
)

type GraphKey struct {
	position  utils.Position
	level     int
	direction utils.Direction
}

type Day17Graph = utils.Graph[GraphKey]

func makeGridIntoGraph(grid utils.Grid) Day17Graph {
	graph := make(Day17Graph)
	for y, row := range grid {
		for x := range row {
			from := utils.Position{X: x, Y: y}
			hasNorth, north := from.North()
			hasEast, east := from.East(len(row) - 1)
			hasWest, west := from.West()
			hasSouth, south := from.South(len(grid) - 1)
			for level := 0; level <= 3; level++ {
				fromEast := GraphKey{position: from, direction: utils.East, level: level}
				fromWest := GraphKey{position: from, direction: utils.West, level: level}
				fromNorth := GraphKey{position: from, direction: utils.North, level: level}
				fromSouth := GraphKey{position: from, direction: utils.South, level: level}
				if hasNorth {
					if level < 3 {
						graph.AddEdge(fromNorth, GraphKey{position: north, direction: utils.North, level: level + 1}, int(grid.ItemAt(north)-'0'))
					}
					graph.AddEdge(fromEast, GraphKey{position: north, direction: utils.North, level: 1}, int(grid.ItemAt(north)-'0'))
					graph.AddEdge(fromWest, GraphKey{position: north, direction: utils.North, level: 1}, int(grid.ItemAt(north)-'0'))
				}
				if hasSouth {
					if level < 3 {
						graph.AddEdge(fromSouth, GraphKey{position: south, direction: utils.South, level: level + 1}, int(grid.ItemAt(south)-'0'))
					}
					graph.AddEdge(fromEast, GraphKey{position: south, direction: utils.South, level: 1}, int(grid.ItemAt(south)-'0'))
					graph.AddEdge(fromWest, GraphKey{position: south, direction: utils.South, level: 1}, int(grid.ItemAt(south)-'0'))
				}
				if hasEast {
					if level < 3 {
						graph.AddEdge(fromEast, GraphKey{position: east, direction: utils.East, level: level + 1}, int(grid.ItemAt(east)-'0'))
					}
					graph.AddEdge(fromSouth, GraphKey{position: east, direction: utils.East, level: 1}, int(grid.ItemAt(east)-'0'))
					graph.AddEdge(fromNorth, GraphKey{position: east, direction: utils.East, level: 1}, int(grid.ItemAt(east)-'0'))
				}
				if hasWest {
					if level < 3 {
						graph.AddEdge(fromWest, GraphKey{position: west, direction: utils.West, level: level + 1}, int(grid.ItemAt(west)-'0'))
					}
					graph.AddEdge(fromSouth, GraphKey{position: west, direction: utils.West, level: 1}, int(grid.ItemAt(west)-'0'))
					graph.AddEdge(fromNorth, GraphKey{position: west, direction: utils.West, level: 1}, int(grid.ItemAt(west)-'0'))
				}
			}
		}
	}
	return graph
}

func renderPath(grid utils.Grid, path []GraphKey) string {
	for _, pos := range path {
		grid.Set(pos.position, '.')
	}
	return grid.String()
}

func part1(fname string) int {
	grid := utils.ReadGridFromFile(fname)
	graph := makeGridIntoGraph(grid)
	directions := []utils.Direction{utils.North, utils.East, utils.West, utils.South}
	froms := make([]GraphKey, 0)
	for _, direction := range directions {
		froms = append(froms, GraphKey{position: utils.Position{X: 0, Y: 0}, level: 0, direction: direction})
	}
	ends := make([]GraphKey, 0)
	for level := 0; level <= 3; level++ {
		for _, direction := range directions {
			ends = append(ends, GraphKey{position: utils.Position{X: len(grid[0]) - 1, Y: len(grid) - 1}, level: level, direction: direction})
		}
	}
	distance, paths := graph.FindDistanceAndPath(froms, ends)
	for path := range paths {
		fmt.Println(paths[path])
		//		fmt.Println(renderPath(grid.Clone(), paths[path]))
	}
	return distance
}

func main() {
	fmt.Println(part1("day17-input-easy2.txt"))
	fmt.Println(part1("day17-input-easy.txt"))
	fmt.Println(part1("day17-input.txt"))
}
