package blackjack

import (
	"bufio"
	"strings"
	"testing"

	deck "github.com/blogscot/card-deck"
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

	t.Run("player scoring", func(t *testing.T) {
		var playerTests = []struct {
			title string
			cards deck.Deck
			score int
		}{
			{"a player can score cards", []deck.Card{four, ten}, 14},
			{"a player has a low ace card", []deck.Card{six, ten, ace}, 17},
			{"a player has a high ace card", []deck.Card{ten, ace}, 21},
		}

		for _, tt := range playerTests {
			player1.hand = tt.cards
			assertEquals(t, player1.score(), tt.score)
		}
	})

	t.Run("dealer scoring", func(t *testing.T) {
		var dealerTests = []struct {
			title string
			cards deck.Deck
			score int
		}{
			{"the dealer can score cards", []deck.Card{ace, king}, 21},
			{"the dealer has a low ace card", []deck.Card{ace, six, jack}, 17},
			{"the dealer has a high ace card", []deck.Card{ace, nine}, 20},
		}

		for _, tt := range dealerTests {
			dealer.hand = tt.cards
			assertEquals(t, dealer.score(), tt.score)
		}
	})
}

func TestWinner(t *testing.T) {

	player1 := Player{name: "TestPlayer"}
	dealer := Dealer{Player: Player{name: "Dealer"}}

	assertWinner := func(want string) {
		t.Helper()

		game := []Participant{&player1, &dealer}
		got := decideWinner(game)

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}

	var winnerTests = []struct {
		title       string
		playerCards deck.Deck
		dealerCards deck.Deck
		winner      string
	}{
		{"player wins by 1", []deck.Card{nine, queen}, []deck.Card{ten, eight}, "TestPlayer"},
		{"dealer wins by 2", []deck.Card{nine, ten}, []deck.Card{ace, ten}, "Dealer"},
		{"game is drawn 1", []deck.Card{nine, ten}, []deck.Card{nine, ten}, gameIsDrawn},
		{"game is drawn 2", []deck.Card{five, ten, four}, []deck.Card{nine, ten}, gameIsDrawn},
	}

	for _, tt := range winnerTests {
		player1.hand = tt.playerCards
		dealer.hand = tt.dealerCards

		assertWinner(tt.winner)
	}
}

func TestPlayerInput(t *testing.T) {

	assertChoice := func(want choice) {
		t.Helper()
		got := playerChoice()

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}

	t.Run("Player hits", func(t *testing.T) {
		sr := strings.NewReader("H\n")
		// Override the standard IO reader
		reader = bufio.NewReader(sr)
		assertChoice(hit)
	})

	t.Run("Player stands", func(t *testing.T) {
		sr := strings.NewReader("s\n")
		reader = bufio.NewReader(sr)

		assertChoice(stand)
	})

	t.Run("Player hits then stands", func(t *testing.T) {
		sr := strings.NewReader("H\nh\nh\ns\n")
		reader = bufio.NewReader(sr)
		assertChoice(hit)
		assertChoice(hit)
		assertChoice(hit)
		assertChoice(stand)
	})
}

func TestParticipants(t *testing.T) {

	player1 = Player{name: "TestPlayer"}
	dealer = Dealer{Player: Player{name: "Dealer"}}

	t.Run("the game aborts if a player goes bust", func(t *testing.T) {
		player1.hand = []deck.Card{jack, four}
		dealer.hand = []deck.Card{ace, six}

		nextCards := deck.Deck([]deck.Card{king, five, six})
		dealer.cards = &nextCards

		// Player hits
		sr := strings.NewReader("h\n")
		reader = bufio.NewReader(sr)

		err := handlePlayer(&player1)
		want := "TestPlayer"

		if err.Error() != want {
			t.Errorf("got %s, want %s", err.Error(), want)
		}
	})

	t.Run("dealer does not go bust on drawing ace with thirteen", func(t *testing.T) {
		player1.hand = []deck.Card{nine, ten}
		dealer.hand = []deck.Card{three, ten}

		nextCards := deck.Deck([]deck.Card{ace, six})
		dealer.cards = &nextCards

		// Player stands
		sr := strings.NewReader("s\n")
		reader = bufio.NewReader(sr)

		ps := []Participant{&player1, &dealer}
		for _, p := range ps {
			_ = play(&p)
		}

		want := 20

		assertEquals(t, dealer.score(), want)

	})

	t.Run("dealer should accept draw if winning is impossible", func(t *testing.T) {
		player1.hand = []deck.Card{ace, ten}
		dealer.hand = []deck.Card{ace, six}

		nextCards := deck.Deck([]deck.Card{four, five, six})
		dealer.cards = &nextCards

		// Player stands
		sr := strings.NewReader("s\n")
		reader = bufio.NewReader(sr)

		ps := []Participant{&player1, &dealer}
		for _, p := range ps {
			_ = play(&p)
		}
		want := 21

		assertEquals(t, dealer.score(), want)
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
	}

	t.Run("Player hits and goes bust", func(t *testing.T) {
		initTest()
		cards := deck.Deck{ten, three, two, five, ten, eight}

		// Player hits and quits game
		sr := strings.NewReader("h\nn\n")
		reader = bufio.NewReader(sr)

		Play(&testingDeck{cards})

		assertEquals(t, player1.score(), 22)
		assertEquals(t, dealer.score(), 8)
	})

	t.Run("Player wins with blackjack", func(t *testing.T) {
		initTest()
		cards := deck.Deck{ten, eight, ace, five, seven, three}

		// Player stands and quits game
		sr := strings.NewReader("s\nn\n")
		reader = bufio.NewReader(sr)

		Play(&testingDeck{cards})

		assertEquals(t, player1.score(), 21)
		assertEquals(t, dealer.score(), 23)
	})
}

func assertEquals(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
