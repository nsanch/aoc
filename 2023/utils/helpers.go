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
