package main

import (
	"fmt"
	"slices"

	"github.com/nsanch/aoc/aoc2023/utils"
)

func getRowsWithoutGalaxies(grid utils.Grid) map[int]bool {
	var ret []int
	for y, row := range grid {
		if !slices.Contains(row, '#') {
			ret = append(ret, y)
		}
	}
	return utils.MakeSetFromSlice(ret)
}

func getColumnsWithoutGalaxies(grid utils.Grid) map[int]bool {
	var ret []int
	for x := 0; x < len(grid[0]); x++ {
		colHasGalaxy := false
		for y := 0; y < len(grid); y++ {
			if grid[y][x] == '#' {
				colHasGalaxy = true
			}
		}
		if !colHasGalaxy {
			ret = append(ret, x)
		}
	}
	return utils.MakeSetFromSlice(ret)
}

func parseFilePart1(fname string) utils.Grid {
	grid := utils.ReadGridFromFile(fname)
	gridWithDoubledRows := make(utils.Grid, 0)
	for _, row := range grid {
		rowHasGalaxy := false
		for _, r := range row {
			if r == '#' {
				rowHasGalaxy = true
				break
			}
		}
		gridWithDoubledRows = append(gridWithDoubledRows, row)
		if !rowHasGalaxy {
			newRow := make([]rune, len(row))
			copy(newRow, row)
			gridWithDoubledRows = append(gridWithDoubledRows, newRow)
		}
	}

	gridWithDoubledColumns := make(utils.Grid, len(gridWithDoubledRows))
	for x := 0; x < len(gridWithDoubledRows[0]); x++ {
		colHasGalaxy := false
		for y := 0; y < len(gridWithDoubledRows); y++ {
			if gridWithDoubledRows[y][x] == '#' {
				colHasGalaxy = true
			}
			gridWithDoubledColumns[y] = append(gridWithDoubledColumns[y], gridWithDoubledRows[y][x])
		}
		if !colHasGalaxy {
			for y := 0; y < len(gridWithDoubledRows); y++ {
				gridWithDoubledColumns[y] = append(gridWithDoubledColumns[y], gridWithDoubledRows[y][x])
			}
		}
	}

	return gridWithDoubledColumns
}

func findGalaxies(grid utils.Grid) []utils.Position {
	var galaxies []utils.Position
	for y, row := range grid {
		for x, r := range row {
			if r == '#' {
				galaxies = append(galaxies, utils.Position{X: x, Y: y})
			}
		}
	}
	return galaxies
}

func part1(fname string) int {
	grid := parseFilePart1(fname)
	galaxies := findGalaxies(grid)
	results := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			manhattanDistance := galaxies[i].ManhattanDistance(galaxies[j])
			results += manhattanDistance
		}
	}
	return results
}

func countDoubledThingsInRange(doubledRows map[int]bool, start int, end int) int {
	if start > end {
		start, end = end, start
	}
	ret := 0
	for i := start + 1; i < end; i++ {
		if doubledRows[i] {
			ret++
		}
	}
	return ret
}

func part2(fname string, scalingFactor int) int {
	grid := utils.ReadGridFromFile(fname)
	galaxies := findGalaxies(grid)
	doubledRows := getRowsWithoutGalaxies(grid)
	doubledColumns := getColumnsWithoutGalaxies(grid)
	results := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			manhattanDistance := galaxies[i].ManhattanDistance(galaxies[j])
			manhattanDistance += countDoubledThingsInRange(doubledRows, galaxies[i].Y, galaxies[j].Y) * (scalingFactor - 1)
			manhattanDistance += countDoubledThingsInRange(doubledColumns, galaxies[i].X, galaxies[j].X) * (scalingFactor - 1)
			results += manhattanDistance
		}
	}
	return results
}

func main() {
	fmt.Println(part1("day11-input-easy.txt"))
	fmt.Println(part1("day11-input.txt"))

	fmt.Println(part2("day11-input-easy.txt", 2))
	fmt.Println(part2("day11-input-easy.txt", 10))
	fmt.Println(part2("day11-input-easy.txt", 100))
	fmt.Println(part2("day11-input.txt", 1000000))
}
