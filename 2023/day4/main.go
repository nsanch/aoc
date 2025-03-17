package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Card struct {
	cardid          int
	numbers         []int
	winning_numbers []int
}

func (card Card) CardValue() int {
	value := 0
	for _, number := range card.numbers {
		for _, winning_number := range card.winning_numbers {
			if number == winning_number {
				if value == 0 {
					value = 1
				} else {
					value = value * 2
				}
			}
		}
	}
	return value
}

func (card Card) NumMatches() int {
	value := 0
	for _, number := range card.numbers {
		for _, winning_number := range card.winning_numbers {
			if number == winning_number {
				value += 1
			}
		}
	}
	return value
}

func convertToInt(s []string) []int {
	var ret []int
	for _, v := range s {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, i)
	}
	return ret
}

func readFile(fname string) []Card {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	var ret []Card
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		cardid_re := regexp.MustCompile(`Card +(\d+): (.*)$`)
		cardid, err := strconv.Atoi(cardid_re.FindStringSubmatch(t)[1])
		if err != nil {
			log.Fatal(err)
		}
		remainder := cardid_re.FindStringSubmatch(t)[2]
		split_by_pipe := strings.Split(remainder, "|")
		card_numbers := strings.Split(split_by_pipe[0], " ")
		winning_numbers := strings.Split(split_by_pipe[1], " ")
		ret = append(ret, Card{cardid: cardid, numbers: convertToInt(card_numbers), winning_numbers: convertToInt(winning_numbers)})
	}
	return ret
}

func part1(fname string) int {
	cards := readFile(fname)
	result := 0
	for _, card := range cards {
		result += card.CardValue()
	}
	fmt.Println(result)
	return result
}

func part2(fname string) int {
	cards := readFile(fname)
	num_copies_of_each_card := make(map[int]int)
	for _, card := range cards {
		num_copies_of_each_card[card.cardid] = 1
	}
	total_number_of_cards := 0
	for _, card := range cards {
		num_copies_of_this_card := num_copies_of_each_card[card.cardid]
		total_number_of_cards += num_copies_of_this_card
		for i := 1; i <= card.NumMatches(); i++ {
			num_copies_of_each_card[card.cardid+i] += num_copies_of_this_card
		}
	}
	fmt.Println(total_number_of_cards)
	return total_number_of_cards
}

func main() {
	part1("day4/day4-input-easy.txt")
	part1("day4/day4-input.txt")

	part2("day4/day4-input-easy.txt")
	part2("day4/day4-input.txt")
}
