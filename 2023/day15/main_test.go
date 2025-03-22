package main

import "testing"

func Test_hash(t *testing.T) {
	tests := []struct {
		args string
		want int
	}{
		{"HASH", 52},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			if got := hash(tt.args); got != tt.want {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
