package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/nsanch/aoc/aoc2023/utils"
)

type Instruction struct {
	direction utils.Direction
	distance  int
	color     string
}

func (i Instruction) String() string {
	return fmt.Sprintf("{Direction: %s, Distance: %v, Color %v}", i.direction.String(), i.distance, i.color)
}

func processInstructions(instructions []Instruction) (utils.Grid, []utils.Position) {
	type GridlessPosition struct {
		x int
		y int
	}
	minX, maxX, minY, maxY := 0, 0, 0, 0
	currPos := GridlessPosition{0, 0}
	for _, instruction := range instructions {
		switch instruction.direction {
		case utils.North:
			currPos.y -= instruction.distance
		case utils.South:
			currPos.y += instruction.distance
		case utils.East:
			currPos.x += instruction.distance
		case utils.West:
			currPos.x -= instruction.distance
		}
		if minX > currPos.x {
			minX = currPos.x
		}
		if maxX < currPos.x {
			maxX = currPos.x
		}
		if minY > currPos.y {
			minY = currPos.y
		}
		if maxY < currPos.y {
			maxY = currPos.y
		}
	}
	gridXBounds := maxX - minX + 1
	gridYBounds := maxY - minY + 1
	grid := make(utils.Grid, gridYBounds)
	for i := range grid {
		grid[i] = make([]rune, gridXBounds)
		copy(grid[i], []rune(strings.Repeat(".", gridXBounds)))
	}
	currPos2 := utils.Position{X: -minX, Y: -minY}
	path := make([]utils.Position, 0)
	path = append(path, currPos2)
	for _, instruction := range instructions {
		for i := 0; i < instruction.distance; i++ {
			switch instruction.direction {
			case utils.North:
				currPos2.Y--
			case utils.South:
				currPos2.Y++
			case utils.East:
				currPos2.X++
			case utils.West:
				currPos2.X--
			}
			grid[currPos2.Y][currPos2.X] = '#'
			path = append(path, currPos2)
		}
	}
	return grid, path
}

func part1(fname string) int {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	instructions := make([]Instruction, 0)
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		re := regexp.MustCompile(`(\w) (\d+) \((#[0-9a-f]{6})\)`)
		matches := re.FindStringSubmatch(t)
		directionStr := matches[1]
		var direction utils.Direction
		switch directionStr {
		case "R":
			direction = utils.East
		case "L":
			direction = utils.West
		case "U":
			direction = utils.North
		case "D":
			direction = utils.South
		}

		distance, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		color := matches[3]
		instruction := Instruction{direction: utils.Direction(direction), distance: distance, color: color}
		instructions = append(instructions, instruction)
	}
	grid, path := processInstructions(instructions)
	if path[len(path)-1] == path[0] {
		path = path[:len(path)-1]
	}
	fmt.Println(grid.String())
	interior := utils.GetInteriorPoints(&grid, path)
	for _, pos := range interior {
		grid.Set(pos, '#')
	}
	fmt.Println(grid.String())
	return grid.Count('#')
}

func main() {
	fmt.Println(part1("day18-input-easy.txt"))
	fmt.Println(part1("day18-input.txt"))
}
