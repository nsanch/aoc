package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/nsanch/aoc/aoc2023/utils"
)

type Row struct {
	springs           string
	sizesOfBrokenSets []int
}

func parseFile(fname string) []Row {
	var ret []Row
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		splitUp := strings.Split(row, " ")
		springs := splitUp[0]
		sizesOfBrokenSets := utils.ConvertStringsToInts(strings.Split(splitUp[1], ","))
		ret = append(ret, Row{springs: springs, sizesOfBrokenSets: sizesOfBrokenSets})
	}
	return ret
}

func checkCompatibility(springs string, sizesOfBrokenSets []int) bool {
	//log.Printf("Checking compatibility for %s with %v\n", springs, sizesOfBrokenSets))

	setsInExample := strings.Split(springs, ".")
	setsInExample = slices.DeleteFunc(setsInExample, func(s string) bool {
		return s == ""
	})
	if len(sizesOfBrokenSets) != len(setsInExample) {
		return false
	}
	for i := 0; i < len(setsInExample); i++ {
		if len(setsInExample[i]) != sizesOfBrokenSets[i] {
			return false
		}
	}
	return true
}

func findNumberOfMatchingSets(springs string, sizesOfBrokenSets []int) int {
	var workingSet []string
	workingSet = make([]string, 0)
	workingSet = append(workingSet, springs)
	numPossibleVersions := 0
	for len(workingSet) > 0 {
		// treat this like a DFS by looking at end of list and appending to end of list.
		curr := workingSet[len(workingSet)-1]
		workingSet = workingSet[:len(workingSet)-1]
		firstQuestionMark := strings.IndexRune(curr, '?')
		if firstQuestionMark == -1 {
			// no more question marks. see if it's compatible with expectation.
			if checkCompatibility(curr, sizesOfBrokenSets) {
				numPossibleVersions++
			}
		} else {
			// replace question mark with . and # and append to list.
			usingDot := curr[:firstQuestionMark] + "." + curr[firstQuestionMark+1:]
			usingHash := curr[:firstQuestionMark] + "#" + curr[firstQuestionMark+1:]
			workingSet = append(workingSet, usingDot)
			workingSet = append(workingSet, usingHash)
		}
	}
	return numPossibleVersions
}

func part1(fname string) int {
	rows := parseFile(fname)
	result := 0
	for _, row := range rows {
		curr := findNumberOfMatchingSets(row.springs, row.sizesOfBrokenSets)
		//fmt.Printf("%d ways to fix %s\n", curr, row.springs)
		result += curr
	}
	return result
}

func main() {
	fmt.Println(part1("day12-input-easy.txt"))
	fmt.Println(part1("day12-input.txt"))
}
