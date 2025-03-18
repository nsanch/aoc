package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nsanch/aoc/aoc2023/utils"
)

type Sequence []int

func (s Sequence) Differentiate() Sequence {
	out := make(Sequence, len(s)-1)
	for i := 1; i < len(s); i++ {
		out[i-1] = s[i] - s[i-1]
	}
	return out
}

func (s Sequence) IsAllZeros() bool {
	for i := range s {
		if s[i] != 0 {
			return false
		}
	}
	return true
}

func (s Sequence) PredictNextValue() int {
	if s.IsAllZeros() {
		return 0
	}
	return s[len(s)-1] + s.Differentiate().PredictNextValue()
}

func (s Sequence) PredictPreviousValue() int {
	if s.IsAllZeros() {
		return 0
	}
	return s[0] - s.Differentiate().PredictPreviousValue()
}

func parseFile(fname string) []Sequence {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var ret []Sequence
	for scanner.Scan() {
		line := scanner.Text()
		sequence := utils.ConvertStringsToInts(strings.Fields(line))
		ret = append(ret, sequence)
	}
	return ret
}

func part1(fname string) int {
	sequences := parseFile(fname)
	result := 0
	for _, seq := range sequences {
		result += seq.PredictNextValue()
	}
	return result
}

func part2(fname string) int {
	sequences := parseFile(fname)
	result := 0
	for _, seq := range sequences {
		result += seq.PredictPreviousValue()
	}
	return result
}

func main() {
	fmt.Println(part1("day9-input-easy.txt"))
	fmt.Println(part1("day9-input.txt"))

	fmt.Println(part2("day9-input-easy.txt"))
	fmt.Println(part2("day9-input.txt"))
}
