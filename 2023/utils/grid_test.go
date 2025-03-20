package utils

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadGridFromFD(t *testing.T) {
	type Test struct {
		in   string
		want Grid
	}
	tests := []Test{{
		in: "123\n456\n789",
		want: Grid{
			{'1', '2', '3'},
			{'4', '5', '6'},
			{'7', '8', '9'},
		}}}

	for _, curr := range tests {
		t.Run(curr.in, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(curr.in))
			ok, got := ReadGridFromFD(scanner)
			if !ok || !got.Equal(curr.want) {
				t.Errorf("ReadGridFromFD(%s) = %v, %v, want `true`, %v", curr.in, ok, got, curr.want)
			}
		})
	}
}
