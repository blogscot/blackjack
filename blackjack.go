// Package blackjack contains functionality to play the card game Blackjack
package blackjack

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blogscot/deck"
)

// Pack provides an interface for testing purposes
type Pack interface {
	Shuffle()
	Cards() *deck.Deck
}

// Player holds a players details
type Player struct {
	name  string
	cards []deck.Card
}

// Dealer holds the dealers details
type Dealer struct {
	Player
	hiddenCard deck.Card
	deal       func() (card deck.Card)
}

type choice int

const (
	hit = iota
	stand
)

const (
	pageBreak    = "=================================\n"
	dealerIsBust = `
++++++++++++++++++++++++++++
The dealer is BUST. You win!
++++++++++++++++++++++++++++`
	playerIsBust = `
++++++++++++
You're BUST!
++++++++++++`
	gameIsDrawn = "Draw"
	letsPlay    = `Let's play BlackJack!

The dealer shuffles the deck thoroughly then starts dealing...
`
)

var (
	reader = bufio.NewReader(os.Stdin)

	dealer  = Dealer{Player: Player{name: "The dealer"}, deal: dealCard}
	player1 = Player{name: "Player1"}
	players = []Participant{&player1, &dealer}
	cards   *deck.Deck
)

// Play plays the game
func Play(pack Pack) {
	fmt.Println(letsPlay)

	cards = pack.Cards()
	cards.Shuffle()

	dealFirstHand()
	showHands(false)

	justPlayersLength := len(players) - 1
	justPlayers := players[:justPlayersLength]
	for _, p := range justPlayers {
		if err := play(&p); err != nil {
			fmt.Printf("%s is BUST!\n", err.Error())
			return
		}
	}

	d := players[justPlayersLength]
	if err := play(&d); err != nil {
		fmt.Printf("The dealer is BUST!\n")
		return
	}

	winner := decideWinner(players)
	if winner == gameIsDrawn {
		fmt.Printf("\nThe game is a draw!")
	} else {
		fmt.Printf("\n%s wins!\n", winner)
	}
}

func dealFirstHand() {
	for n := 1; n <= 2; n++ {
		for _, p := range players {
			c := dealCard()
			p.add(c)
		}
	}
}

func dealCard() (card deck.Card) {
	card = (*cards)[0]
	*cards = (*cards)[1:]
	return
}

func (p *Player) getName() string {
	return p.name
}

func (p *Player) add(c deck.Card) {
	fmt.Printf("%s receives %s\n", p.name, c)
	p.cards = append(p.cards, c)
}

func (d *Dealer) add(c deck.Card) {
	if d.hiddenCard == (deck.Card{}) {
		fmt.Printf("The dealer places a card face down.\n")
		d.hiddenCard = c
	} else {
		fmt.Printf("and receives %s\n", c)
		d.cards = append(d.cards, c)
	}
}

func (p Player) score() (total int) {
	for _, c := range p.cards {
		total += scoreCard(c)
	}
	if p.hasAce() && total <= 11 {
		total += 10
	}
	return
}

func (d Dealer) score() (total int) {
	for _, c := range d.cards {
		total += scoreCard(c)
	}
	total += scoreCard(d.hiddenCard)
	if d.hasAce() && total <= 11 {
		total += 10
	}
	return
}

func (p Player) hasAce() bool {
	for _, c := range p.cards {
		if c.Value == deck.Ace {
			return true
		}
	}
	return false
}

func (d Dealer) hasAce() bool {
	return d.Player.hasAce() || d.hiddenCard.Value == deck.Ace
}

func scoreCard(c deck.Card) (score int) {
	switch c.Value {
	case deck.Ace:
		score = 1
	case deck.Two:
		score = 2
	case deck.Three:
		score = 3
	case deck.Four:
		score = 4
	case deck.Five:
		score = 5
	case deck.Six:
		score = 6
	case deck.Seven:
		score = 7
	case deck.Eight:
		score = 8
	case deck.Nine:
		score = 9
	case deck.Ten:
		score = 10
	case deck.Jack:
		score = 10
	case deck.Queen:
		score = 10
	case deck.King:
		score = 10
	default:
		log.Fatalf("invalid card: %s of %s", c.Value, c.Suit)
	}
	return
}

func playerChoice() choice {
	var (
		text    string
		isValid = false
	)

	for !isValid {
		fmt.Print(pageBreak)
		fmt.Print("Do you want to (H)it or (S)tand? ")
		text, _ = reader.ReadString('\n')
		text = strings.TrimRight(strings.ToLower(text), "\n")
		if text == "h" || text == "s" {
			isValid = true
		}
	}
	if text == "h" {
		return hit
	}
	return stand
}

// Standard wraps the deck to facilitate testing
type Standard struct {
	deck.Deck
}

// Shuffle shuffles the standard deck of cards
func (p *Standard) Shuffle() {
	p.Deck.Shuffle()
}

// Cards unwraps the deck
func (p Standard) Cards() *deck.Deck {
	return &p.Deck
}

// Deck returns a brand new card deck.
func Deck() Standard {
	return Standard{deck.New()}
}
