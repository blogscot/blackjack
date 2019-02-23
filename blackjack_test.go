package main

import (
	"testing"

	"github.com/blogscot/deck"
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
		card2 := deck.Card{Suit: deck.Hearts, Value: deck.Ten}
		player1.cards = []deck.Card{card1, card2}
		wanted := 14
		got := player1.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("a player's has a low ace card", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Spades, Value: deck.Six}
		card2 := deck.Card{Suit: deck.Spades, Value: deck.Ten}
		card3 := deck.Card{Suit: deck.Diamonds, Value: deck.Ace}
		player1.cards = []deck.Card{card1, card2, card3}
		wanted := 17
		got := player1.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("a player's has a high ace card", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Spades, Value: deck.Ace}
		card2 := deck.Card{Suit: deck.Spades, Value: deck.Ten}
		player1.cards = []deck.Card{card1, card2}
		wanted := 21
		got := player1.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("the dealer can score cards", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Clubs, Value: deck.Ace}
		card2 := deck.Card{Suit: deck.Hearts, Value: deck.King}
		dealer.faceDown = card1
		dealer.cards = []deck.Card{card2}

		wanted := 21
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("the dealer has a low ace card", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Spades, Value: deck.Six}
		card2 := deck.Card{Suit: deck.Spades, Value: deck.Ten}
		card3 := deck.Card{Suit: deck.Diamonds, Value: deck.Ace}
		dealer.cards = []deck.Card{card1, card2}
		dealer.faceDown = card3
		wanted := 17
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("the dealer has a high ace card", func(t *testing.T) {
		card1 := deck.Card{Suit: deck.Clubs, Value: deck.Nine}
		card2 := deck.Card{Suit: deck.Diamonds, Value: deck.Ace}
		dealer.cards = []deck.Card{card1}
		dealer.faceDown = card2
		wanted := 20
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

}
