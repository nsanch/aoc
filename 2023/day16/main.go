package main

import (
	"fmt"
	"slices"

	"github.com/nsanch/aoc/aoc2023/utils"
)

type VisitedKey struct {
	Direction utils.Direction
	Position  utils.Position
}

type VisitedMap map[VisitedKey]bool

func (m *VisitedMap) DidVisit(pos utils.Position, direction utils.Direction) bool {
	_, ok := (*m)[VisitedKey{Direction: direction, Position: pos}]
	return ok
}

func (m *VisitedMap) Visit(pos utils.Position, direction utils.Direction) {
	(*m)[VisitedKey{Direction: direction, Position: pos}] = true
}

func simulateLight(grid utils.Grid, energyGrid utils.Grid, lightPos utils.Position, lightDirection utils.Direction, visited *VisitedMap) {
	if visited.DidVisit(lightPos, lightDirection) {
		return
	}
	(*visited).Visit(lightPos, lightDirection)

	//fmt.Printf("lightPos: %v, lightDirection: %v. item: %c\n", lightPos, lightDirection, grid.ItemAt(lightPos))
	energyGrid.Set(lightPos, '#')
	switch grid.ItemAt(lightPos) {
	case '/':
		switch lightDirection {
		case utils.North:
			lightDirection = utils.East
		case utils.East:
			lightDirection = utils.North
		case utils.South:
			lightDirection = utils.West
		case utils.West:
			lightDirection = utils.South
		}
	case '\\':
		switch lightDirection {
		case utils.North:
			lightDirection = utils.West
		case utils.West:
			lightDirection = utils.North
		case utils.South:
			lightDirection = utils.East
		case utils.East:
			lightDirection = utils.South
		}
	case '|':
		if lightDirection == utils.East || lightDirection == utils.West {
			simulateLight(grid, energyGrid, lightPos, utils.South, visited)
			simulateLight(grid, energyGrid, lightPos, utils.North, visited)
			return
		}
	case '-':
		if lightDirection == utils.South || lightDirection == utils.North {
			simulateLight(grid, energyGrid, lightPos, utils.West, visited)
			simulateLight(grid, energyGrid, lightPos, utils.East, visited)
			return
		}
	}
	hasNext, next := lightPos.FollowDirection(lightDirection, len(grid[0])-1, len(grid)-1)
	if hasNext {
		simulateLight(grid, energyGrid, next, lightDirection, visited)
	}
}

func ScoreGrid(grid utils.Grid) int {
	score := 0
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == '#' {
				score++
			}
		}
	}
	return score
}

func SimulateLightForGrid(grid utils.Grid, startingPos utils.Position, startingDirection utils.Direction) int {
	energyGrid := make(utils.Grid, len(grid))
	for i := range energyGrid {
		energyGrid[i] = slices.Repeat([]rune{'.'}, len(grid[i]))
	}
	visitedMap := make(VisitedMap)
	simulateLight(grid, energyGrid, startingPos, startingDirection, &visitedMap)
	return ScoreGrid(energyGrid)
}

func part1(fname string) int {
	grid := utils.ReadGridFromFile(fname)
	return SimulateLightForGrid(grid, utils.Position{Y: 0, X: 0}, utils.East)
}

func part2(fname string) int {
	grid := utils.ReadGridFromFile(fname)
	bestScore := 0
	type Start struct {
		pos       utils.Position
		direction utils.Direction
	}
	toSimulate := make([]Start, 0)
	for y := range grid {
		toSimulate = append(toSimulate, Start{utils.Position{Y: y, X: 0}, utils.East})
		toSimulate = append(toSimulate, Start{utils.Position{Y: y, X: len(grid[y]) - 1}, utils.West})
	}
	for x := range grid[0] {
		toSimulate = append(toSimulate, Start{utils.Position{Y: 0, X: x}, utils.South})
		toSimulate = append(toSimulate, Start{utils.Position{Y: len(grid) - 1, X: x}, utils.North})
	}
	for _, start := range toSimulate {
		score := SimulateLightForGrid(grid, start.pos, start.direction)
		if score > bestScore {
			bestScore = score
		}
	}
	return bestScore
}

func main() {
	fmt.Println(part1("day16-input-easy.txt"))
	fmt.Println(part1("day16-input.txt"))

	fmt.Println(part2("day16-input-easy.txt"))
	fmt.Println(part2("day16-input.txt"))
}
