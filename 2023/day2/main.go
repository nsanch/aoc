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

type Batch struct {
	green_balls int
	red_balls   int
	blue_balls  int
}

type Game struct {
	gameid  int
	batches []Batch
}

func parseFile(fname string) []Game {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	games := make([]Game, 0)
	for scanner.Scan() {
		t := scanner.Text()
		gameid_re := regexp.MustCompile(`Game (\d+): (.*)$`)
		gameid_matches := gameid_re.FindStringSubmatch(t)
		gameid, err := strconv.Atoi(gameid_matches[1])
		if err != nil {
			log.Fatal(err)
		}
		game := Game{gameid: gameid, batches: nil}
		remainder := gameid_matches[2]
		split_up := strings.Split(remainder, ";")
		batches := make([]Batch, len(split_up))
		for pos, one_batch := range split_up {
			b := &batches[pos]
			for _, s2 := range strings.Split(one_batch, ",") {
				num_and_color := strings.Split(strings.TrimSpace(s2), " ")
				num, err := strconv.Atoi(num_and_color[0])
				if err != nil {
					log.Fatal(err)
				}
				color := num_and_color[1]
				switch color {
				case "red":
					b.red_balls = num
				case "green":
					b.green_balls = num
				case "blue":
					b.blue_balls = num
				}
			}
		}
		game.batches = batches
		games = append(games, game)
		//matches := re.FindAllString(t, -1)
		//u1, _ := strconv.Atoi(matches[0] + matches[len(matches)-1])
		//total += u1
	}
	//fmt.Println(games)
	return games
}

func isGamePossible(game Game, starting_red int, starting_green int, starting_blue int) bool {
	for _, batch := range game.batches {
		if batch.red_balls > starting_red || batch.blue_balls > starting_blue || batch.green_balls > starting_green {
			return false
		}
	}
	return true
}

func part1(fname string) int {
	games := parseFile(fname)
	result := 0
	for _, game := range games {
		if isGamePossible(game, 12, 13, 14) {
			result += game.gameid
		}
	}
	return result
}

func part2(fname string) int {
	games := parseFile(fname)
	result := 0
	for _, game := range games {
		var max_red, max_blue, max_green int
		for _, batch := range game.batches {
			max_red = max(max_red, batch.red_balls)
			max_blue = max(max_blue, batch.blue_balls)
			max_green = max(max_green, batch.green_balls)
		}
		game_power := max_red * max_blue * max_green
		result += game_power
	}
	return result
}

func main() {
	fmt.Println(part1("day2/day2-input-easy.txt"))
	fmt.Println(part1("day2/day2-input.txt"))

	fmt.Println(part2("day2/day2-input-easy.txt"))
	fmt.Println(part2("day2/day2-input.txt"))
}
