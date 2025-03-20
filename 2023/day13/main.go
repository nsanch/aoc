package main

import (
	"fmt"
	"slices"

	"github.com/nsanch/aoc/aoc2023/utils"
)

func findNumberOfReflectingRows(grid utils.Grid) int {
	bestSoFar := 0
	for y := 1; y < len(grid); y++ {
		// we have a reflecting row, back up and see how many more we can find.
		gotToAnEdge := true
		for toBackUp := 0; y-toBackUp-1 >= 0 && y+toBackUp < len(grid); toBackUp++ {
			if !slices.Equal(grid[y-toBackUp-1], grid[y+toBackUp]) {
				gotToAnEdge = false
			}
		}
		if gotToAnEdge {
			bestSoFar = y
		}
	}
	return bestSoFar
}

func part1(fname string) int {
	grids := utils.ReadGridsFromFile(fname)
	result := 0
	for _, grid := range grids {
		result += 100 * findNumberOfReflectingRows(grid)
		transposed := grid.Transpose()
		result += findNumberOfReflectingRows(transposed)
	}
	return result
}

func main() {
	fmt.Println(part1("day13-input-easy.txt"))
	fmt.Println(part1("day13-input.txt"))
}
