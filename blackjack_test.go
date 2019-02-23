package main

import (
	"testing"

	"github.com/blogscot/deck"
)

var (
	ten = deck.Card{Suit: deck.Hearts, Value: deck.Ten}
	ace = deck.Card{Suit: deck.Spades, Value: deck.Ace}
)

func TestScoring(t *testing.T) {

	player1 := Player{
		name: "TestPlayer",
	}

	dealer := Dealer{
		Player: Player{
			name: "Dealer",
		},
	}

	t.Run("a player can score cards", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Hearts, Value: deck.Four}
		player1.cards = []deck.Card{card1, ten}
		wanted := 14
		got := player1.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("a player has a low ace card", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Spades, Value: deck.Six}
		player1.cards = []deck.Card{card1, ten, ace}
		wanted := 17
		got := player1.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("a player has a high ace card", func(t *testing.T) {
		player1.cards = []deck.Card{ace, ten}
		wanted := 21
		got := player1.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("the dealer can score cards", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Hearts, Value: deck.King}
		dealer.faceDown = ace
		dealer.cards = []deck.Card{card1}

		wanted := 21
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("the dealer has a low ace card", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Spades, Value: deck.Six}
		dealer.cards = []deck.Card{card1, ten}
		dealer.faceDown = ace
		wanted := 17
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("the dealer has a high ace card", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Clubs, Value: deck.Nine}
		dealer.cards = []deck.Card{card1}
		dealer.faceDown = ace
		wanted := 20
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

}
