package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/nsanch/aoc/aoc2023/utils"
)

type Race struct {
	duration int
	distance int
}

func parseFile(name string, ignoreSpaces bool) []Race {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	timeLine := scanner.Text()

	scanner.Scan()
	distanceLine := scanner.Text()

	if ignoreSpaces {
		timeLine = strings.ReplaceAll(timeLine, " ", "")
		distanceLine = strings.ReplaceAll(distanceLine, " ", "")
	}
	numbersRE := regexp.MustCompile(`\d+`)
	times := utils.ConvertStringsToInts(numbersRE.FindAllString(timeLine, -1))
	distances := utils.ConvertStringsToInts(numbersRE.FindAllString(distanceLine, -1))
	if len(times) != len(distances) {
		log.Fatalf("Mismatched number of times and distances. %d vs %d", len(times), len(distances))
	}
	var races []Race
	for i := 0; i < len(times); i++ {
		races = append(races, Race{duration: times[i], distance: distances[i]})
	}

	return races
}

func findNumberOfWinningSolutions(race Race) int {
	wins := 0
	for pressLength := 0; pressLength < race.duration; pressLength++ {
		// if we press for `pressLength` seconds, we will have traveled `race.pressLength * (race.duration - pressLength)` distance.
		// if that is more than race.distance, we have won.
		distance := pressLength * (race.duration - pressLength)
		if distance > race.distance {
			wins++
		}
	}
	return wins
}

func part1(fname string) int {
	races := parseFile(fname, false)
	result := 1
	for _, race := range races {
		result = result * findNumberOfWinningSolutions(race)
	}
	return result
}

func part2(fname string) int {
	races := parseFile(fname, true)
	return findNumberOfWinningSolutions(races[0])
}

func main() {
	fmt.Println(part1("day6-input-easy.txt"))
	fmt.Println(part1("day6-input.txt"))

	fmt.Println(part2("day6-input-easy.txt"))
	fmt.Println(part2("day6-input.txt"))
}
