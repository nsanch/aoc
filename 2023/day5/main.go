package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/nsanch/aoc/aoc2023/utils"
)

type InputMapRange struct {
	destinationRangeStart int
	sourceRangeStart      int
	sourceRangeEnd        int
}

type InputMap struct {
	sourceName      string
	destinationName string
	definedRanges   []InputMapRange
}

func (inputMap InputMap) mapInputToOutput(input int) int {
	for _, definedRange := range inputMap.definedRanges {
		if input >= definedRange.sourceRangeStart && input < definedRange.sourceRangeEnd {
			return definedRange.destinationRangeStart + (input - definedRange.sourceRangeStart)
		}
	}
	// if there's no explicit mapping, it maps to itself.
	return input
}

func (inputMap InputMap) mapOutputToInput(output int) int {
	for _, definedRange := range inputMap.definedRanges {
		destRangeEnd := definedRange.destinationRangeStart + (definedRange.sourceRangeEnd - definedRange.sourceRangeStart)
		if output >= definedRange.destinationRangeStart && output < destRangeEnd {
			return definedRange.sourceRangeStart + (output - definedRange.destinationRangeStart)
		}
	}
	// if there's no explicit mapping, it maps to itself.
	return output
}

type CategoryAndLocation struct {
	category string
	location int
}

type Almanac struct {
	desiredSeeds []CategoryAndLocation
	inputMaps    map[string][]InputMap
}

func (almanac Almanac) getAllOutputsForInput(input int, kind string) []CategoryAndLocation {
	relevantMaps := almanac.inputMaps[kind]
	outputs := make([]CategoryAndLocation, len(relevantMaps))
	for _, inputMap := range relevantMaps {
		outputs = append(outputs, CategoryAndLocation{
			category: inputMap.destinationName,
			location: inputMap.mapInputToOutput(input)})
	}
	return outputs
}

func (almanac Almanac) getAllInputsForOutput(output int, outputKind string) []CategoryAndLocation {
	var ret []CategoryAndLocation
	for _, inputMaps := range almanac.inputMaps {
		for _, inputMap := range inputMaps {
			if inputMap.destinationName == outputKind {
				ret = append(ret, CategoryAndLocation{
					category: inputMap.sourceName,
					location: inputMap.mapOutputToInput(output)})
			}
		}
	}
	return ret
}

func parseFile(name string) *Almanac {
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var almanac *Almanac = new(Almanac)
	almanac.inputMaps = make(map[string][]InputMap)
	scanner.Scan()
	seedsLine := scanner.Text()
	seeds_re := regexp.MustCompile(`\d+`)
	desiredSeeds := utils.ConvertStringsToInts(seeds_re.FindAllString(seedsLine, -1))
	for _, seed := range desiredSeeds {
		almanac.desiredSeeds = append(almanac.desiredSeeds, CategoryAndLocation{
			category: "seed", location: seed})
	}
	scanner.Scan() // consume a newline.
	for scanner.Scan() {
		mapLine := scanner.Text()
		mapLineRE := regexp.MustCompile(`(.*)-to-(.*) map`)
		sourceMapName := mapLineRE.FindStringSubmatch(mapLine)[1]
		destMapName := mapLineRE.FindStringSubmatch(mapLine)[2]
		var definedRanges []InputMapRange
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}
			lineNumbers := utils.ConvertStringsToInts(strings.Split(strings.TrimSpace(line), " "))
			definedRanges = append(definedRanges, InputMapRange{
				destinationRangeStart: lineNumbers[0],
				sourceRangeStart:      lineNumbers[1],
				sourceRangeEnd:        lineNumbers[1] + lineNumbers[2]})
		}
		almanac.inputMaps[sourceMapName] = append(almanac.inputMaps[sourceMapName], InputMap{
			sourceName:      sourceMapName,
			destinationName: destMapName,
			definedRanges:   definedRanges})
	}
	return almanac
}

func navigateMapsToLocations(almanac *Almanac, startingPoints []CategoryAndLocation) []int {
	knownData := make([]CategoryAndLocation, len(startingPoints))
	copy(knownData, startingPoints)
	var results []int
	for len(knownData) > 0 {
		// take the first element of knownData and map it to everything we can
		// add that to the end of knownData
		// repeat.
		current := knownData[0]
		knownData = knownData[1:]
		outputs := almanac.getAllOutputsForInput(current.location, current.category)
		for _, output := range outputs {
			knownData = append(knownData, output)
			if output.category == "location" {
				results = append(results, output.location)
			}
		}
	}
	return results
}

func navigateBackToSeedFromLocation(almanac *Almanac, startingPoints []CategoryAndLocation) []CategoryAndLocation {
	knownData := make([]CategoryAndLocation, len(startingPoints))
	copy(knownData, startingPoints)
	var results []CategoryAndLocation
	for len(knownData) > 0 {
		// take the first element of knownData and map it to everything we can
		// add that to the end of knownData
		// repeat.
		current := knownData[0]
		knownData = knownData[1:]
		inputs := almanac.getAllInputsForOutput(current.location, current.category)
		for _, input := range inputs {
			if input.category == "seed" {
				results = append(results, input)
			} else {
				knownData = append(knownData, input)
			}
		}
	}
	return results
}

func part1(fname string) int {
	almanac := parseFile(fname)
	locations := navigateMapsToLocations(almanac, almanac.desiredSeeds)
	fmt.Println(locations)
	return slices.Min(locations)
}

/*
	func navigateFromOneSeedLocationToBest(almanac *Almanac, seed CategoryAndLocation) int {
		locations := navigateMapsToLocations(almanac, []CategoryAndLocation{seed})
		if len(locations) == 0 {
			return 0
		}
		return slices.Min(locations)
	}
*/
func part2(fname string) int {
	almanac := parseFile(fname)
	for currLocation := 1; true; currLocation++ {
		start := CategoryAndLocation{category: "location", location: currLocation}
		results := navigateBackToSeedFromLocation(almanac, []CategoryAndLocation{start})
		for _, result := range results {
			for pairStart := 0; pairStart < len(almanac.desiredSeeds); pairStart += 2 {
				start := almanac.desiredSeeds[pairStart].location
				end := start + almanac.desiredSeeds[pairStart+1].location
				if result.location >= start &&
					result.location < end {
					fmt.Printf("Started from %d and got to seed %d\n", currLocation, result.location)
					return currLocation
				}
			}
		}
		if currLocation%1000000 == 0 {
			fmt.Printf("Continuing on after %d\n", currLocation)
		}
	}
	// shouldn't happen since the above is an infinite loop, but go wants a return here.
	return 0
}

func main() {
	fmt.Println(part1("day5-input-easy.txt"))
	fmt.Println(part1("day5-input.txt"))

	fmt.Println(part2("day5-input-easy.txt"))
	fmt.Println(part2("day5-input.txt"))
}
