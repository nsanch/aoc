package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/nsanch/aoc/aoc2023/utils"
)

func findAllPositionsOfReflectingRows(grid utils.Grid) []int {
	var ret []int
	for y := 1; y < len(grid); y++ {
		// we have a reflecting row, back up and see how many more we can find.
		gotToAnEdge := true
		for toBackUp := 0; y-toBackUp-1 >= 0 && y+toBackUp < len(grid); toBackUp++ {
			if !slices.Equal(grid[y-toBackUp-1], grid[y+toBackUp]) {
				gotToAnEdge = false
			}
		}
		if gotToAnEdge {
			ret = append(ret, y)
		}
	}
	return ret
}

func part1(fname string) int {
	grids := utils.ReadGridsFromFile(fname)
	result := 0
	for _, grid := range grids {
		rowPositions := findAllPositionsOfReflectingRows(grid)
		colPositions := findAllPositionsOfReflectingRows(grid.Transpose())
		if len(rowPositions) > 0 {
			result += 100 * slices.Max(rowPositions)
		}
		if len(colPositions) > 0 {
			result += slices.Max(colPositions)
		}
	}
	return result
}

func fixExactlyOneDifference(a, b []rune) (bool, []rune, int) {
	if len(a) != len(b) {
		log.Fatal("can't supply different length arguments to hasExactlyOneDifference")
	}
	replacementA := make([]rune, len(a))
	copy(replacementA, a)
	diffCount := 0
	replacedPos := -1
	for i := range a {
		if a[i] != b[i] {
			diffCount++
			replacementA[i] = b[i]
			replacedPos = i
		}
	}
	if diffCount != 1 {
		return false, nil, -1
	} else {
		return true, replacementA, replacedPos
	}
}

func findPositionOfReflectingRowsWithASmudge(grid utils.Grid) int {
	oldReflections := utils.MakeSetFromSlice(findAllPositionsOfReflectingRows(grid))
	allNewIndexes := make([]int, 0)
	for y := range grid {
		for y2 := y + 1; y2 < len(grid); y2++ {
			hasDiff, replacement, _ := fixExactlyOneDifference(grid[y], grid[y2])
			if hasDiff {
				clonedGrid := grid.Clone()
				clonedGrid[y] = replacement
				newReflections := findAllPositionsOfReflectingRows(clonedGrid)
				for newIdx := range newReflections {
					if !oldReflections[newIdx] {
						//fmt.Printf("found a new reflection at row %d after replacing (y=%d,x=%d)\n", newIdx, y, posX)
						allNewIndexes = append(allNewIndexes, newIdx)
					}
				}

				clonedGrid[y2] = replacement
				newReflections = findAllPositionsOfReflectingRows(clonedGrid)
				for _, newIdx := range newReflections {
					if !oldReflections[newIdx] {
						//fmt.Printf("found a new reflection at row %d after replacing (y=%d,x=%d)\n", newIdx, y, posX)
						allNewIndexes = append(allNewIndexes, newIdx)
					}
				}
			}
		}
	}
	if len(allNewIndexes) == 0 {
		return 0
	}
	return slices.Max(allNewIndexes)
}

func part2(fname string) int {
	grids := utils.ReadGridsFromFile(fname)
	result := 0
	for _, grid := range grids {
		result += 100 * findPositionOfReflectingRowsWithASmudge(grid)
		transposed := grid.Transpose()
		result += findPositionOfReflectingRowsWithASmudge(transposed)
	}
	return result
}

func main() {
	fmt.Println(part1("day13-input-easy.txt"))
	fmt.Println(part1("day13-input.txt"))

	fmt.Println(part2("day13-input-easy.txt"))
	fmt.Println(part2("day13-input.txt"))
}
