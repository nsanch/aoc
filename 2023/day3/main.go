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

func ReadFile(fname string) Schematic {
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
	s := ReadFile(fname)
	part_numbers := findPartNumbers(s)
	result := 0
	for _, num := range part_numbers {
		result += num
	}
	return result
}

// PART 2

type PartAndNeighbor struct {
	neighborPosition Position
	partNumber       int
}

func getNeighbors(s Schematic, partNumber int, startPos Position, numCharacters int) []PartAndNeighbor {
	var ret []PartAndNeighbor
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
					if s.grid[y][x] == '*' {
						ret = append(ret, PartAndNeighbor{neighborPosition: Position{x: x, y: y}, partNumber: partNumber})
					}
				}
			}
		}
	}
	return ret
}

func getPartNumberNeighbors(s Schematic, number_str []rune, startPos Position) []PartAndNeighbor {
	partNum, err := strconv.Atoi(string(number_str))
	if err != nil {
		log.Fatal(err)
	}
	return getNeighbors(s, partNum, startPos, len(number_str))
}

func findSymbols(s Schematic) []PartAndNeighbor {
	var ret []PartAndNeighbor
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
				neighbors := getPartNumberNeighbors(s, number_str, startPos)
				ret = append(ret, neighbors...)
				number_str = nil
			}
		}
		if len(number_str) > 0 {
			neighbors := getPartNumberNeighbors(s, number_str, startPos)
			ret = append(ret, neighbors...)
		}
	}
	return ret
}

func Part2(fname string) int {
	s := ReadFile(fname)
	symbol_positions := findSymbols(s)
	var m map[Position][]int = make(map[Position][]int)
	result := 0

	for _, pos := range symbol_positions {
		partNumbers, in_map := m[pos.neighborPosition]
		if !in_map {
			m[pos.neighborPosition] = []int{pos.partNumber}
		} else if len(partNumbers) == 1 {
			gearRatio := partNumbers[0] * pos.partNumber
			result += gearRatio
		}
	}
	return result
}

func main() {
	fmt.Println(part1("day3/day3-input-easy.txt"))
	fmt.Println(part1("day3/day3-input.txt"))

	fmt.Println(Part2("day3/day3-input-easy.txt"))
	fmt.Println(Part2("day3/day3-input.txt"))
}
