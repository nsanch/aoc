package main

import (
	"bufio"
	"cmp"
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

func (r Row) Compare(other Row) int {
	return cmp.Or(
		cmp.Compare(r.springs, other.springs),
		slices.Compare(r.sizesOfBrokenSets, other.sizesOfBrokenSets))
}

func (r Row) Repeat() Row {
	return Row{
		springs:           strings.Join(slices.Repeat([]string{r.springs}, 5), "?"),
		sizesOfBrokenSets: slices.Repeat(r.sizesOfBrokenSets, 5)}
	//log.Print(r.springs)
	//log.Print(r.sizesOfBrokenSets)
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

func GetBrokenSetSizes(springs string) []int {
	//log.Printf("Checking compatibility for %s with %v\n", springs, sizesOfBrokenSets))

	setsInExample := strings.Split(springs, ".")
	setsInExample = slices.DeleteFunc(setsInExample, func(s string) bool {
		return s == ""
	})
	var ret []int
	for _, str := range setsInExample {
		ret = append(ret, len(str))
	}
	return ret
}

func checkCompatibility(springs string, sizesOfBrokenSets []int) bool {
	//log.Printf("Checking compatibility for %s with %v\n", springs, sizesOfBrokenSets))
	setsInGivenString := GetBrokenSetSizes(springs)
	return slices.Equal(setsInGivenString, sizesOfBrokenSets)
}

type MatchingSetsCacheValue struct {
	sizesOfBrokenSets   []int
	numPossibleVersions int
}

type MatchingSetsCache map[string][]MatchingSetsCacheValue

func (cache MatchingSetsCache) Get(springs string, sizesOfBrokenSets []int) int {
	for _, value := range cache[springs] {
		if slices.Compare(value.sizesOfBrokenSets, sizesOfBrokenSets) == 0 {
			return value.numPossibleVersions
		}
	}
	return -1
}

func (cache MatchingSetsCache) Set(springs string, sizesOfBrokenSets []int, numPossibleVersions int) int {
	currValue := cache.Get(springs, sizesOfBrokenSets)
	if currValue != -1 {
		if currValue != numPossibleVersions {
			log.Fatalf("for %s and %v, expected %d but got %d", springs, sizesOfBrokenSets, currValue, numPossibleVersions)
		}
		return currValue
	}

	if cache[springs] == nil {
		cache[springs] = make([]MatchingSetsCacheValue, 0)
	}
	cache[springs] = append(cache[springs], MatchingSetsCacheValue{sizesOfBrokenSets: sizesOfBrokenSets, numPossibleVersions: numPossibleVersions})
	return numPossibleVersions
}

type ConsumptionCacheKey struct {
	springs        string
	firstBrokenSet int
}

type ConsumptionCache map[ConsumptionCacheKey][]string

func consumePrefixOfStringSatisfyingBrokenSet(springs string, firstBrokenSet int, cache *ConsumptionCache) []string {
	cacheKey := ConsumptionCacheKey{springs: springs, firstBrokenSet: firstBrokenSet}
	cachedValue, ok := (*cache)[cacheKey]
	if ok {
		return cachedValue
	}

	//log.Printf("checking if %s has a prefix of %d\n", springs, firstBrokenSet)

	if len(springs) < firstBrokenSet {
		(*cache)[cacheKey] = nil
		return nil
	}

	firstQuestionMark := strings.IndexRune(springs, '?')
	// check the N-character string starting with the first ? or hash.
	firstHash := strings.IndexRune(springs, '#')
	if firstQuestionMark == -1 && firstHash == -1 {
		(*cache)[cacheKey] = nil
		return nil
	}

	if firstQuestionMark == -1 || (firstHash != -1 && firstHash < firstQuestionMark) {
		if firstHash+firstBrokenSet > len(springs) {
			(*cache)[cacheKey] = nil
			return nil
		}
		stringStartingAtHash := springs[firstHash : firstHash+firstBrokenSet]
		for _, ch := range stringStartingAtHash {
			if ch == '.' {
				(*cache)[cacheKey] = nil
				return nil
			}
		}
		// great, that worked as long as next char isn't a # too
		if len(springs) > (firstHash+firstBrokenSet) && springs[firstHash+firstBrokenSet] == '#' {
			(*cache)[cacheKey] = nil
			return nil
		}
		// set this to a dot in case it's a question mark to ensure it's not a hash.
		solution := ""
		if len(springs) >= firstHash+firstBrokenSet+1 {
			solution = "." + springs[firstHash+firstBrokenSet+1:]
		} else {
			solution = "."
		}
		// success. recurse.
		ret := []string{solution}
		//log.Printf("success with %s %d\n", springs[firstHash:firstHash+firstBrokenSet], firstBrokenSet)
		(*cache)[cacheKey] = ret
		return ret
	} else {
		usingDot := "." + springs[firstQuestionMark+1:]
		usingHash := "#" + springs[firstQuestionMark+1:]
		remainderUsingDot := consumePrefixOfStringSatisfyingBrokenSet(usingDot, firstBrokenSet, cache)
		remainderUsingHash := consumePrefixOfStringSatisfyingBrokenSet(usingHash, firstBrokenSet, cache)
		//log.Printf("%d solutions with %s %d\n", len(remainderUsingDot), usingDot, firstBrokenSet)
		//log.Printf("%d solutions with %s %d\n", len(remainderUsingHash), usingHash, firstBrokenSet)
		ret := slices.Concat(remainderUsingDot, remainderUsingHash)
		(*cache)[cacheKey] = ret
		return ret
	}
}

func findNumberOfMatchingSetsRecursive(springs string, sizesOfBrokenSets []int, cache *MatchingSetsCache, cache2 *ConsumptionCache) int {
	cachedValue := cache.Get(springs, sizesOfBrokenSets)
	if cachedValue != -1 {
		return cachedValue
	}

	firstQuestionMark := strings.IndexRune(springs, '?')
	if firstQuestionMark == -1 {
		// no more question marks. see if it's compatible with expectation.
		if checkCompatibility(springs, sizesOfBrokenSets) {
			return cache.Set(springs, sizesOfBrokenSets, 1)
		} else {
			return cache.Set(springs, sizesOfBrokenSets, 0)
		}
	}

	// consume the first element of `sizesOfBrokenSets` and recurse with the remainder.
	if len(sizesOfBrokenSets) == 0 {
		if !strings.ContainsRune(springs, '#') {
			// only one way can work -- set all question marks to dots.
			return cache.Set(springs, sizesOfBrokenSets, 1)
		} else {
			return cache.Set(springs, sizesOfBrokenSets, 0)
		}
	}

	possibleValidRemainders := consumePrefixOfStringSatisfyingBrokenSet(springs, sizesOfBrokenSets[0], cache2)
	if len(possibleValidRemainders) == 0 {
		return cache.Set(springs, sizesOfBrokenSets, 0)
	}

	numPossibleVersions := 0
	for _, remainder := range possibleValidRemainders {
		// recurse with the remainder.
		numPossibleVersions += findNumberOfMatchingSetsRecursive(remainder, sizesOfBrokenSets[1:], cache, cache2)
	}
	//log.Printf("%s with %v has %d possible versions\n", springs, sizesOfBrokenSets, numPossibleVersions)
	return cache.Set(springs, sizesOfBrokenSets, numPossibleVersions)
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

func part2(fname string) int {
	rows := parseFile(fname)
	unfoldedRows := make([]Row, 0)
	for _, row := range rows {
		unfoldedRows = append(unfoldedRows, row.Repeat())
	}
	rows = unfoldedRows
	result := 0
	cache := make(MatchingSetsCache)
	cache2 := make(ConsumptionCache)
	for _, row := range rows {
		curr := findNumberOfMatchingSetsRecursive(row.springs, row.sizesOfBrokenSets, &cache, &cache2)
		//fmt.Printf("%d ways to fix %s / %v\n", curr, row.springs, row.sizesOfBrokenSets) //, strings.Join(answers, ", "))
		result += curr
	}
	return result
}

func main() {
	fmt.Println(part1("day12-input-easy.txt"))
	fmt.Println(part1("day12-input.txt"))

	fmt.Println(part2("day12-input-easy.txt"))
	fmt.Println(part2("day12-input.txt"))
}
