package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Schematic struct {
	grid [][]rune
}

type Position struct {
	x int
	y int
}

func readFile(fname string) Schematic {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	grid := make([][]rune, 0)
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		row := make([]rune, 0)
		for _, r := range t {
			row = append(row, r)
		}
		grid = append(grid, row)
	}
	return Schematic{grid: grid}
}

func hasNeighbor(s Schematic, startPos Position, numCharacters int) bool {
	// look at any character neighboring the part number, including diagonals,
	// so it's like a box surrounding the number. ignore the number itself and any '.' characters.
	for x := startPos.x - 1; x <= startPos.x+numCharacters; x++ {
		for y := startPos.y - 1; y <= startPos.y+1; y++ {
			if y >= 0 && y < len(s.grid) && x >= 0 && x < len(s.grid[y]) {
				if x >= startPos.x && x < startPos.x+numCharacters && y == startPos.y {
					continue
				}
				if s.grid[y][x] != '.' && !unicode.IsDigit(s.grid[y][x]) {
					//fmt.Printf("found neighbor %c at %d,%d relative to %d,%d+%d\n", s.grid[y][x], x, y, startPos.x, startPos.y, numCharacters)
					return true
				}
			}
		}
	}
	return false
}

func checkPartNumber(s Schematic, number_str []rune, startPos Position) int {
	number, err := strconv.Atoi(string(number_str))
	if err != nil {
		log.Fatal(err)
	}
	if hasNeighbor(s, startPos, len(number_str)) {
		//fmt.Println("found part number", number)
		return number
	}
	return 0
}

func findPartNumbers(s Schematic) []int {
	part_numbers := make([]int, len(s.grid))
	for y, row := range s.grid {
		//fmt.Printf("ROW %d\n", y)
		var number_str []rune
		var startPos Position
		for x, r := range row {
			if unicode.IsDigit(r) {
				if number_str == nil {
					startPos = Position{x: x, y: y}
				}
				number_str = append(number_str, r)
			} else if len(number_str) > 0 {
				number := checkPartNumber(s, number_str, startPos)
				if number > 0 {
					part_numbers = append(part_numbers, number)
				}
				number_str = nil
			}
		}
		if len(number_str) > 0 {
			number := checkPartNumber(s, number_str, startPos)
			if number > 0 {
				part_numbers = append(part_numbers, number)
			}
		}
	}
	return part_numbers
}

func part1(fname string) int {
	s := readFile(fname)
	part_numbers := findPartNumbers(s)
	result := 0
	for _, num := range part_numbers {
		result += num
	}
	return result
}

func main() {
	fmt.Println(part1("day3/day3-input-easy.txt"))
	fmt.Println(part1("day3/day3-input.txt"))
}
