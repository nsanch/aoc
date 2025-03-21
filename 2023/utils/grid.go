package utils

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strings"
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type Position struct {
	X int
	Y int
}

func (p Position) East(maxXValue int) (bool, Position) {
	if p.X < maxXValue {
		return true, Position{X: p.X + 1, Y: p.Y}
	}
	return false, Position{}
}

func (p Position) West() (bool, Position) {
	if p.X > 0 {
		return true, Position{X: p.X - 1, Y: p.Y}
	}
	return false, Position{}
}

func (p Position) North() (bool, Position) {
	if p.Y > 0 {
		return true, Position{X: p.X, Y: p.Y - 1}
	}
	return false, Position{}
}

func (p Position) South(maxYValue int) (bool, Position) {
	if p.Y < maxYValue {
		return true, Position{X: p.X, Y: p.Y + 1}
	}
	return false, Position{}
}

func (p Position) FollowDirection(d Direction, maxXValue int, maxYValue int) (bool, Position) {
	switch d {
	case North:
		return p.North()
	case East:
		return p.East(maxXValue)
	case West:
		return p.West()
	case South:
		return p.South(maxYValue)
	}
	log.Fatal("invalid direction", d)
	return false, Position{}
}

func (p Position) ManhattanDistance(other Position) int {
	return Abs(p.X-other.X) + Abs(p.Y-other.Y)
}

type Grid [][]rune

func (grid Grid) ItemAt(pos Position) rune {
	return grid[pos.Y][pos.X]
}

func (grid Grid) Set(pos Position, value rune) {
	grid[pos.Y][pos.X] = value
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
