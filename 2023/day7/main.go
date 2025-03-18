package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Card int

const (
	Joker Card = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func makeCardFromRune(r rune, useJokerRule bool) Card {
	switch r {
	case '2':
		return Two
	case '3':
		return Three
	case '4':
		return Four
	case '5':
		return Five
	case '6':
		return Six
	case '7':
		return Seven
	case '8':
		return Eight
	case '9':
		return Nine
	case 'T':
		return Ten
	case 'J':
		if useJokerRule {
			return Joker
		} else {
			return Jack
		}
	case 'Q':
		return Queen
	case 'K':
		return King
	case 'A':
		return Ace
	}
	log.Fatalf("Invalid card %c", r)
	return -1
}

func (c Card) String() string {
	switch c {
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "T"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	}
	log.Fatalf("Invalid card %d", c)
	return ""
}

type Hand struct {
	// in part1, J means Jack, but in part2, J means Joker
	useJokerRule         bool
	cardsInOriginalOrder []Card
	cards                []Card
	bid                  int
}

type HandKind int

const (
	HighCard HandKind = iota
	OnePair
	TwoPairs
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

func (kind HandKind) String() string {
	switch kind {
	case HighCard:
		return "HighCard"
	case OnePair:
		return "OnePair"
	case TwoPairs:
		return "TwoPairs"
	case ThreeOfAKind:
		return "ThreeOfAKind"
	case FullHouse:
		return "FullHouse"
	case FourOfAKind:
		return "FourOfAKind"
	case FiveOfAKind:
		return "FiveOfAKind"
	}
	log.Fatalf("Invalid hand kind %d", kind)
	return ""
}

func (hand Hand) String() string {
	var sb strings.Builder
	sb.WriteString(hand.GetKind().String())
	sb.WriteString(" ")
	for _, card := range hand.CardsInOrderOfRelevance() {
		sb.WriteString(card.String())
	}
	return sb.String()
}

func (hand Hand) GroupCards() map[Card]int {
	groupedCards := make(map[Card]int)
	for _, card := range hand.cards {
		groupedCards[card]++
	}

	return groupedCards
}

func (hand Hand) CardsInOrderOfRelevance() []Card {
	groupedCards := hand.GroupCards()
	switch hand.GetKind() {
	case HighCard:
		return hand.cards
	case OnePair:
		results := make([]Card, 0)
		for card, count := range groupedCards {
			if count == 2 {
				results = append(results, card)
			}
		}
		// since hand.cards is sorted, we can just append the rest of the cards
		for _, card := range hand.cards {
			if card != results[0] {
				results = append(results, card)
			}
		}
		return results
	case TwoPairs:
		var foundPair1 = false
		var pair1 Card
		var pair2 Card
		var nonPair Card
		for card, count := range groupedCards {
			if count == 2 {
				if foundPair1 {
					if pair1 < card {
						pair2 = pair1
						pair1 = card
					} else {
						pair2 = card
					}
				} else {
					pair1 = card
					foundPair1 = true
				}
			} else {
				nonPair = card
			}
		}
		return []Card{pair1, pair2, nonPair}
	case ThreeOfAKind:
		results := make([]Card, 0)
		for card, count := range groupedCards {
			if count == 3 {
				results = append(results, card)
			}
		}
		// since hand.cards is sorted, we can just append the rest of the cards
		for _, card := range hand.cards {
			if card != results[0] {
				results = append(results, card)
			}
		}
		return results
	case FullHouse:
		var three Card
		var two Card
		for card, count := range groupedCards {
			if count == 3 {
				three = card
			} else if count == 2 {
				two = card
			}
		}
		return []Card{three, two}
	case FourOfAKind:
		var four Card
		var otherCard Card
		for card, count := range groupedCards {
			if count == 4 {
				four = card
			} else {
				otherCard = card
			}
		}
		return []Card{four, otherCard}
	case FiveOfAKind:
		return hand.cards[:1]
	}
	log.Fatalf("Invalid hand kind %d", hand.GetKind())
	return nil
}

func (hand Hand) GetKind() HandKind {
	groupedCards := hand.GroupCards()
	switch {
	case len(groupedCards) == 5:
		return HighCard
	case len(groupedCards) == 4:
		return OnePair
	case len(groupedCards) == 3:
		for _, count := range groupedCards {
			if count == 2 {
				//fmt.Println(groupedCards, "TwoPairs")
				return TwoPairs
			}
		}
		//fmt.Println(groupedCards, "ThreeOfAKind")
		return ThreeOfAKind
	case len(groupedCards) == 2:
		// if it's a full-house the counts will be 3&2.
		// If it's a four-of-a-kind the counts will be 4&1
		countForCard0 := groupedCards[hand.cards[0]]
		if countForCard0 == 3 || countForCard0 == 2 {
			return FullHouse
		}
		return FourOfAKind
	case len(groupedCards) == 1:
		return FiveOfAKind
	}
	log.Fatal("Invalid hand")
	return 0
}

func (hand Hand) GetBestVersionGivenJokers() Hand {
	potentialHands := make([]Hand, 0)
	for replacementCard := Two; replacementCard <= Ace; replacementCard++ {
		replacementCards := make([]Card, 0)
		for i := range hand.cards {
			if hand.cards[i] != Joker {
				replacementCards = append(replacementCards, hand.cards[i])
			} else {
				replacementCards = append(replacementCards, replacementCard)
			}
		}
		potentialHands = append(potentialHands, Hand{
			cards:                replacementCards,
			bid:                  hand.bid,
			cardsInOriginalOrder: hand.cardsInOriginalOrder,
			useJokerRule:         false}) // false to avoid infinite recursion
	}
	slices.SortFunc(potentialHands, CompareHands)
	return potentialHands[len(potentialHands)-1]
}

func CompareHands(h1 Hand, h2 Hand) int {
	if h1.useJokerRule {
		h1 = h1.GetBestVersionGivenJokers()
		h2 = h2.GetBestVersionGivenJokers()
	}
	// for some reason we compare based on the original order of the cards not based on the
	// second/third/fourth-ranked card in the hand. this wasted a bunch of time.
	return cmp.Or(
		cmp.Compare(h1.GetKind(), h2.GetKind()),
		slices.Compare(h1.cardsInOriginalOrder, h2.cardsInOriginalOrder))
	//		slices.Compare(h1.CardsInOrderOfRelevance(), h2.CardsInOrderOfRelevance()))
}

func parseFile(fname string, useJokerRule bool) []Hand {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hands []Hand
	for scanner.Scan() {
		splitUp := strings.Split(scanner.Text(), " ")
		bidInt, err := strconv.Atoi(splitUp[1])
		if err != nil {
			log.Fatal(err)
		}
		var cards []Card
		for _, c := range splitUp[0] {
			cards = append(cards, makeCardFromRune(c, useJokerRule))
		}
		origCards := make([]Card, len(cards))
		copy(origCards, cards)
		slices.Sort(cards)
		slices.Reverse(cards)
		hand := Hand{cards: cards, cardsInOriginalOrder: origCards, bid: bidInt, useJokerRule: useJokerRule}
		hands = append(hands, hand)
	}
	return hands
}

func part1(fname string) int {
	hands := parseFile(fname, false)
	slices.SortFunc(hands, CompareHands)
	//fmt.Println(hands)
	result := 0
	for rank, hand := range hands {
		//fmt.Println(hand.String())
		result += hand.bid * (rank + 1)
	}
	return result
}

func part2(fname string) int {
	hands := parseFile(fname, true)
	slices.SortFunc(hands, CompareHands)
	//fmt.Println(hands)
	result := 0
	for rank, hand := range hands {
		//fmt.Println(hand.String())
		result += hand.bid * (rank + 1)
	}
	return result
}

func main() {
	fmt.Println(part1("day7-input-easy.txt"))
	fmt.Println(part1("day7-input.txt"))

	fmt.Println(part2("day7-input-easy.txt"))
	fmt.Println(part2("day7-input.txt"))
}
