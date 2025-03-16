package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func part1(fname string) int {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	total := 0
	for scanner.Scan() {
		t := scanner.Text()
		re := regexp.MustCompile(`\d`)
		matches := re.FindAllString(t, -1)
		u1, _ := strconv.Atoi(matches[0] + matches[len(matches)-1])
		total += u1
	}
	return total
}

func convertMatch(s string) int {
	switch {
	case s == "one":
		return 1
	case s == "two":
		return 2
	case s == "three":
		return 3
	case s == "four":
		return 4
	case s == "five":
		return 5
	case s == "six":
		return 6
	case s == "seven":
		return 7
	case s == "eight":
		return 8
	case s == "nine":
		return 9
	default:
		ret, _ := strconv.Atoi(s)
		return ret
	}
}

func part2(fname string) int {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	total := 0
	for scanner.Scan() {
		t := scanner.Text()
		re := regexp.MustCompile(`(one|two|three|four|five|six|seven|eight|nine|\d)`)
		first := convertMatch(re.FindStringSubmatch(t)[1])
		last := -1
		for i := len(t) - 1; i >= 0; i-- {
			finalMatch := re.FindStringSubmatch(t[i:])
			if finalMatch != nil {
				last = convertMatch(finalMatch[1])
				break
			}
		}
		u1 := (10 * first) + last
		total += u1
	}
	return total
}

func main() {
	fmt.Println(part1("day1/day1-input-easy.txt"))
	fmt.Println(part1("day1/day1-input.txt"))

	fmt.Println(part2("day1/day1-input-easy2.txt"))
	fmt.Println(part2("day1/day1-input.txt"))
}
