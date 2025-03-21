package main

import (
	"fmt"
	"slices"

	"github.com/nsanch/aoc/aoc2023/utils"
)

func getIterationOrder(grid utils.Grid, d utils.Direction) []utils.Position {
	iterationOrder := make([]utils.Position, 0)
	switch d {
	case utils.North:
		fallthrough
	case utils.West:
		for y := range grid {
			for x := range grid[y] {
				iterationOrder = append(iterationOrder, utils.Position{Y: y, X: x})
			}
		}

	case utils.South:
		for y := len(grid) - 1; y >= 0; y-- {
			for x := range grid[y] {
				iterationOrder = append(iterationOrder, utils.Position{Y: y, X: x})
			}
		}
	case utils.East:
		for y := range grid {
			for x := len(grid[y]) - 1; x >= 0; x-- {
				iterationOrder = append(iterationOrder, utils.Position{Y: y, X: x})
			}
		}
	}
	return iterationOrder
}

func tiltGrid(grid utils.Grid, d utils.Direction) bool {
	iterationOrder := getIterationOrder(grid, d)

	tiltedAnything := false
	for _, curr := range iterationOrder {
		if grid.ItemAt(curr) == 'O' {
			hasNext, next := curr.FollowDirection(d, len(grid[0])-1, len(grid)-1)
			for hasNext && grid.ItemAt(next) == '.' {
				tiltedAnything = true
				//fmt.Printf("moving from %d,%d to %d,%d\n", y, x, tiltedY, x)
				grid.Set(curr, '.')
				curr = next
				hasNext, next = curr.FollowDirection(d, len(grid[0])-1, len(grid)-1)
			}
			grid.Set(curr, 'O')
		}
	}
	return tiltedAnything
}

func ScoreGrid(grid utils.Grid) int {
	result := 0
	for y, row := range grid {
		for _, r := range row {
			if r == 'O' {
				result += len(grid) - y
			}
		}
	}
	return result
}

func part1(fname string) int {
	grid := utils.ReadGridFromFile(fname)
	tiltGrid(grid, utils.North)
	result := ScoreGrid(grid)
	return result
}

func spinGrid(grid utils.Grid) {
	tiltGrid(grid, utils.North)
	tiltGrid(grid, utils.West)
	tiltGrid(grid, utils.South)
	tiltGrid(grid, utils.East)
}

func lookForRepeatingPattern(input []int) ([]int, int) {
	startOfPattern := slices.Index(input[1:], input[0]) + 1
	if startOfPattern == 0 {
		pattern, offset := lookForRepeatingPattern(input[1:])
		return pattern, offset + 1
	}
	for i := 0; i < startOfPattern; i++ {
		if input[i] != input[startOfPattern+i] {
			pattern, offset := lookForRepeatingPattern(input[1:])
			return pattern, offset + 1
		}
	}
	// confirm the pattern repeats for full input.
	for i := range input {
		if input[i] != input[(i%startOfPattern)] {
			pattern, offset := lookForRepeatingPattern(input[1:])
			return pattern, offset + 1
		}
	}
	return input[:startOfPattern], 0
}

func part2(fname string) int {
	grid := utils.ReadGridFromFile(fname)
	//fmt.Println(grid.String())
	fmt.Println()

	spinGrid(grid)

	allPriorScores := make([]int, 0)
	allPriorScores = append(allPriorScores, ScoreGrid(grid))

	for i := range 1000 {
		_ = i
		spinGrid(grid)
		allPriorScores = append(allPriorScores, ScoreGrid(grid))
	}
	pattern, offset := lookForRepeatingPattern(allPriorScores)
	result := pattern[(1000000000-offset-1)%len(pattern)]
	fmt.Printf("Found repeating pattern of %v starting at offset %d\n", pattern, offset)
	return result
}

func main() {
	fmt.Println(part1("day14-input-easy.txt"))
	fmt.Println(part1("day14-input.txt"))

	fmt.Println(part2("day14-input-easy.txt"))
	fmt.Println(part2("day14-input.txt"))
}
