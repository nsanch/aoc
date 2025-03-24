package utils

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strings"
)

type Grid [][]rune

func (grid Grid) ItemAt(pos Position) rune {
	return grid[pos.Y][pos.X]
}

func (grid Grid) Set(pos Position, value rune) {
	grid[pos.Y][pos.X] = value
}

func (grid Grid) Count(r rune) int {
	count := 0
	for _, row := range grid {
		for _, v := range row {
			if v == r {
				count++
			}
		}
	}
	return count
}

func (grid Grid) Transpose() Grid {
	transposed := make(Grid, len(grid[0]))
	for i := range transposed {
		transposed[i] = make([]rune, len(grid))
	}
	for y, row := range grid {
		for x, r := range row {
			transposed[x][y] = r
		}
	}
	return transposed
}

func (grid Grid) Equal(other Grid) bool {
	if len(grid) != len(other) {
		return false
	}
	for y, row := range grid {
		if !slices.Equal(row, other[y]) {
			return false
		}
	}
	return true
}

func (grid Grid) Clone() Grid {
	clone := make(Grid, len(grid))
	for y, row := range grid {
		clone[y] = make([]rune, len(row))
		copy(clone[y], row)
	}
	return clone
}

func (grid Grid) String() string {
	var sb strings.Builder
	for _, row := range grid {
		sb.WriteString(string(row))
		sb.WriteString("\n")
	}
	return sb.String()
}

func ReadGridFromFile(fname string) Grid {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	_, grid := ReadGridFromFD(scanner)
	return grid
}

func ReadGridFromFD(scanner *bufio.Scanner) (bool, Grid) {
	grid := make(Grid, 0)
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		if t == "" {
			break
		}
		row := make([]rune, 0)
		for _, r := range t {
			row = append(row, r)
		}
		grid = append(grid, row)
	}
	return len(grid) > 0, grid
}

func ReadGridsFromFile(fname string) []Grid {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	grids := make([]Grid, 0)
	for {
		ok, grid := ReadGridFromFD(scanner)
		if !ok {
			break
		}
		grids = append(grids, grid)
	}
	return grids
}
