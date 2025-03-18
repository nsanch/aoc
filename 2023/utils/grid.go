package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
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

type Grid [][]rune

func (grid Grid) ItemAt(pos Position) rune {
	return grid[pos.Y][pos.X]
}

func ReadGridFromFile(fname string) Grid {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	grid := make(Grid, 0)
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		row := make([]rune, 0)
		for _, r := range t {
			row = append(row, r)
		}
		grid = append(grid, row)
	}
	return grid
}
