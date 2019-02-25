package blackjack

import (
	"bufio"
	"strings"
	"testing"

	"github.com/blogscot/deck"
)

var (
	ace   = deck.Card{Suit: deck.Spades, Value: deck.Ace}
	two   = deck.Card{Suit: deck.Spades, Value: deck.Two}
	three = deck.Card{Suit: deck.Spades, Value: deck.Three}
	four  = deck.Card{Suit: deck.Spades, Value: deck.Four}
	five  = deck.Card{Suit: deck.Spades, Value: deck.Five}
	six   = deck.Card{Suit: deck.Spades, Value: deck.Six}
	seven = deck.Card{Suit: deck.Spades, Value: deck.Seven}
	eight = deck.Card{Suit: deck.Spades, Value: deck.Eight}
	nine  = deck.Card{Suit: deck.Spades, Value: deck.Nine}
	ten   = deck.Card{Suit: deck.Spades, Value: deck.Ten}
	jack  = deck.Card{Suit: deck.Spades, Value: deck.Jack}
	queen = deck.Card{Suit: deck.Spades, Value: deck.Queen}
	king  = deck.Card{Suit: deck.Spades, Value: deck.King}
)

func TestScoring(t *testing.T) {

	player1 := Player{name: "TestPlayer"}
	dealer := Dealer{Player: Player{name: "Dealer"}}

	t.Run("a player can score cards", func(t *testing.T) {
		player1.cards = []deck.Card{four, ten}
		wanted := 14
		got := player1.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("a player has a low ace card", func(t *testing.T) {
		player1.cards = []deck.Card{six, ten, ace}
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
		dealer.hiddenCard = ace
		dealer.cards = []deck.Card{king}

		wanted := 21
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("the dealer has a low ace card", func(t *testing.T) {
		dealer.cards = []deck.Card{six, jack}
		dealer.hiddenCard = ace
		wanted := 17
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("the dealer has a high ace card", func(t *testing.T) {
		dealer.cards = []deck.Card{nine}
		dealer.hiddenCard = ace
		wanted := 20
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})
}

func TestWinner(t *testing.T) {

	player1 := Player{name: "TestPlayer"}
	dealer := Dealer{Player: Player{name: "Dealer"}}

	t.Run("player wins", func(t *testing.T) {
		player1.cards = []deck.Card{nine, queen}
		dealer.cards = []deck.Card{eight}
		dealer.hiddenCard = ten
		game := []Participant{&player1, &dealer}
		wanted := "TestPlayer"

		got := decideWinner(game)

		if got != wanted {
			t.Errorf("got %q, wanted %q", got, wanted)
		}
	})

	t.Run("dealer wins", func(t *testing.T) {
		player1.cards = []deck.Card{nine, ten}
		dealer.cards = []deck.Card{ten}
		dealer.hiddenCard = ace
		game := []Participant{&player1, &dealer}
		wanted := "Dealer"

		got := decideWinner(game)

		if got != wanted {
			t.Errorf("got %q, wanted %q", got, wanted)
		}
	})

	t.Run("game is drawn", func(t *testing.T) {
		player1.cards = []deck.Card{nine, ten}
		dealer.cards = []deck.Card{ten}
		dealer.hiddenCard = nine
		game := []Participant{&player1, &dealer}
		wanted := "Draw"

		got := decideWinner(game)

		if got != wanted {
			t.Errorf("got %q, wanted %q", got, wanted)
		}
	})

}

func TestPlayerInput(t *testing.T) {

	assertChoice := func(t *testing.T, want choice) {
		t.Helper()
		got := playerChoice()

		if got != want {
			t.Errorf("got %q, wanted %q", got, want)
		}
	}

	t.Run("Player hits", func(t *testing.T) {
		sr := strings.NewReader("H\n")
		// Override the standard IO reader
		reader = bufio.NewReader(sr)
		assertChoice(t, hit)
	})

	t.Run("Player stands", func(t *testing.T) {
		sr := strings.NewReader("s\n")
		reader = bufio.NewReader(sr)

		assertChoice(t, stand)
	})

	t.Run("Player hits then stands", func(t *testing.T) {
		sr := strings.NewReader("H\nh\nh\ns\n")
		reader = bufio.NewReader(sr)
		assertChoice(t, hit)
		assertChoice(t, hit)
		assertChoice(t, hit)
		assertChoice(t, stand)
	})
}

func TestDealer(t *testing.T) {

	player1 = Player{name: "TestPlayer"}
	dealer = Dealer{Player: Player{name: "Dealer"}}
	dealer.deal = dealCard

	t.Run("dealer does not go bust on drawing ace with thirteen", func(t *testing.T) {
		player1.cards = []deck.Card{nine, ten}
		dealer.cards = []deck.Card{ten}
		dealer.hiddenCard = three

		nextCards := deck.Deck([]deck.Card{ace, six})
		cards = &nextCards

		// Player stands
		sr := strings.NewReader("s\n")
		reader = bufio.NewReader(sr)

		ps := []Participant{&player1, &dealer}
		for _, p := range ps {
			play(&p)
		}
		wanted := 20
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})

	t.Run("dealer should accept draw if winning is impossible", func(t *testing.T) {
		player1.cards = []deck.Card{ace, ten}
		dealer.cards = []deck.Card{six}
		dealer.hiddenCard = ace

		nextCards := deck.Deck([]deck.Card{four, five, six})
		cards = &nextCards

		// Player stands
		sr := strings.NewReader("s\n")
		reader = bufio.NewReader(sr)

		ps := []Participant{&player1, &dealer}
		for _, p := range ps {
			play(&p)
		}
		wanted := 21
		got := dealer.score()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	})
}

type testingDeck struct {
	deck.Deck
}

func (t *testingDeck) Shuffle() {
	// No shuffling here
}

func (t *testingDeck) Cards() *deck.Deck {
	return &t.Deck
}

func TestGame(t *testing.T) {

	initTest := func() {
		player1 = Player{name: "TestPlayer"}
		dealer = Dealer{Player: Player{name: "Dealer"}}
		dealer.deal = dealCard
	}

	assertEquals := func(t *testing.T, got, wanted int) {
		t.Helper()

		if got != wanted {
			t.Errorf("got %d, wanted %d", got, wanted)
		}
	}

	t.Run("Player hits and goes bust", func(t *testing.T) {
		initTest()
		cards := deck.Deck{ten, three, two, five, ten}

		// Player hits
		sr := strings.NewReader("h\n")
		reader = bufio.NewReader(sr)

		Play(&testingDeck{cards})

		assertEquals(t, player1.score(), 22)
		assertEquals(t, dealer.score(), 8)
	})

	t.Run("Player wins with blackjack", func(t *testing.T) {
		initTest()
		cards := deck.Deck{ten, eight, ace, five, seven, three}

		// Player stands
		sr := strings.NewReader("s\n")
		reader = bufio.NewReader(sr)

		Play(&testingDeck{cards})

		assertEquals(t, player1.score(), 21)
		assertEquals(t, dealer.score(), 23)
	})
}
