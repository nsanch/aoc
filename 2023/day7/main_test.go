package main

import (
	"testing"
)

func TestGetKind(t *testing.T) {
	h := Hand{cards: []Card{Two, Two, Two, Two, Two}, bid: 1}
	if h.GetKind() != FiveOfAKind {
		t.Errorf("Expected FiveOfAKind, got %v", h.GetKind())
	}

	h = Hand{cards: []Card{King, King, Four, Three, Two}, bid: 1}
	if h.GetKind() != OnePair {
		t.Errorf("Expected OnePair, got %v", h.GetKind())
	}

	h = Hand{cards: []Card{King, King, Seven, Seven, Six}, bid: 1}
	if h.GetKind() != TwoPairs {
		t.Errorf("Expected TwoPairs, got %v", h.GetKind())
	}

	h = Hand{cards: []Card{King, Jack, Jack, Ten, Ten}, bid: 1}
	if h.GetKind() != TwoPairs {
		t.Errorf("Expected TwoPairs, got %v", h.GetKind())
	}
}
