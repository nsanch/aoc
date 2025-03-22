package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func hash(s string) int {
	h := 0
	for _, ch := range s {
		h += int(ch)
		h *= 17
		h %= 256
	}
	return h
}

func parseFile(fname string) []string {
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	line := scanner.Text()
	return strings.Split(line, ",")
}

func part1(fname string) int {
	result := 0
	for _, s := range parseFile(fname) {
		result += hash(s)
	}
	return result
}

type BoxSlot struct {
	lens        string
	focalLength int
}

type Box struct {
	slots []BoxSlot
}

func (b *Box) Score() int {
	ret := 0
	for i := range b.slots {
		ret += (i + 1) * b.slots[i].focalLength
	}
	return ret
}

func (b *Box) Add(lens string, focalLength int) {
	for pos, slot := range b.slots {
		if slot.lens == lens {
			b.slots[pos].focalLength = focalLength
			return
		}
	}
	b.slots = append(b.slots, BoxSlot{lens, focalLength})
}

func (b *Box) Remove(lens string) {
	for pos, slot := range b.slots {
		if slot.lens == lens {
			b.slots = slices.Concat(b.slots[:pos], b.slots[pos+1:])
			return
		}
	}
}

type Boxes struct {
	boxes []Box
}

func (b *Boxes) Add(lens string, focalLength int) {
	if len(b.boxes) == 0 {
		b.boxes = make([]Box, 256)
	}
	b.boxes[hash(lens)].Add(lens, focalLength)
}

func (b *Boxes) Remove(lens string) {
	b.boxes[hash(lens)].Remove(lens)
}

func (b *Boxes) Score() int {
	ret := 0
	for pos, box := range b.boxes {
		ret += (pos + 1) * box.Score()
	}
	return ret
}

func part2(fname string) int {
	boxes := new(Boxes)
	lines := parseFile(fname)
	for _, line := range lines {
		re := regexp.MustCompile(`(\w+)(-|=)(\d*)`)
		matches := re.FindStringSubmatch(line)
		lens := matches[1]
		focalLength, _ := strconv.Atoi(matches[3])
		operation := matches[2]
		switch operation {
		case "=":
			boxes.Add(lens, focalLength)
		case "-":
			boxes.Remove(lens)
		}
	}
	return boxes.Score()
}

func main() {
	fmt.Println(part1("day15-input-easy.txt"))
	fmt.Println(part1("day15-input.txt"))

	fmt.Println(part2("day15-input-easy.txt"))
	fmt.Println(part2("day15-input.txt"))
}
