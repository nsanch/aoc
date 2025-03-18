package utils

import (
	"log"
	"strconv"
	"strings"
)

func ConvertStringsToInts(s []string) []int {
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

func GreatestCommonDivisor(a int, b int) int {
	if b > a {
		return GreatestCommonDivisor(b, a)
	}
	// https://en.wikipedia.org/wiki/Greatest_common_divisor#Euclidean_algorithm
	if b == 0 {
		return a
	}
	return GreatestCommonDivisor(b, a%b)
}

func LeastCommonMultiple(a int, b int) int {
	return a * (b / GreatestCommonDivisor(a, b))
}
