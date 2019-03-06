// Package blackjack contains functionality to play the card game Blackjack
package blackjack

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	deck "github.com/blogscot/card-deck"
)

// Pack provides an interface for testing purposes
type Pack interface {
	Shuffle()
	Cards() *deck.Deck
}

// Player holds a players details
type Player struct {
	name string
	hand []deck.Card
}

// Dealer holds the dealers details
type Dealer struct {
	Player
	hiddenCard deck.Card
	cards      *deck.Deck
}

type choice int

const (
	hit = iota
	stand
)

const (
	pageBreak   = "=================================\n"
	gameIsDrawn = "Draw"
	letsPlay    = `Let's play BlackJack!

The dealer shuffles the deck thoroughly.`
	dealNewHand = `The dealer deals a new hand.
`
)

var (
	reader = bufio.NewReader(os.Stdin)

	dealer       = Dealer{Player: Player{name: "The dealer"}}
	player1      = Player{name: "Player1"}
	participants = []Participant{&player1, &dealer}
)

// Play plays the game
func Play(pack Pack) {
	pack.Shuffle()
	gameIntro()

	dealer.cards = pack.Cards()

	for {
		fmt.Println(dealNewHand)
		playHand()

		if playAgain() == "n" {
			break
		}

		fmt.Println()
		clearTable()
	}
	fmt.Println("Thanks for playing. Bye!")
}

func gameIntro() {
	fmt.Println(letsPlay)
}

func playAgain() (text string) {
	fmt.Print("\nPlay again? (Y/n)? ")
	text, _ = reader.ReadString('\n')
	text = strings.TrimRight(strings.ToLower(text), "\n")
	return
}

func playHand() {
	dealer.dealFirstHand()
	showHands(false)

	p := Participant(&player1)
	if err := play(&p); err != nil {
		fmt.Printf("%s is BUST!\n", err.Error())
		return
	}

	d := Participant(&dealer)
	if err := play(&d); err != nil {
		fmt.Printf("The dealer is BUST!\n")
		return
	}

	winner := decideWinner(participants)
	if winner == gameIsDrawn {
		fmt.Printf("\nThe game is a draw!")
	} else {
		fmt.Printf("\n%s wins!\n", winner)
	}
}

func (d Dealer) dealFirstHand() {
	for n := 1; n <= 2; n++ {
		for _, p := range participants {
			c := d.dealCard()
			p.add(c)
		}
	}
}

func (d Dealer) dealCard() (card deck.Card) {
	if len(*d.cards) == 0 {
		fmt.Println("The dealer opens a new deck, and shuffles the cards")
		*d.cards = deck.New()
		d.cards.Shuffle()
	}
	card = (*d.cards)[0]
	*d.cards = (*d.cards)[1:]
	return
}

func (p *Player) getName() string {
	return p.name
}

func (p *Player) add(c deck.Card) {
	p.hand = append(p.hand, c)
}

func (d *Dealer) add(c deck.Card) {
	if d.hiddenCard == (deck.Card{}) {
		d.hiddenCard = c
	} else {
		d.hand = append(d.hand, c)
	}
}

func clearTable() {
	player1.hand = []deck.Card{}
	dealer.hand = []deck.Card{}
	dealer.hiddenCard = deck.Card{}
}

func (p Player) score() (total int) {
	for _, c := range p.hand {
		total += scoreCard(c)
	}
	if p.hasAce() && total <= 11 {
		total += 10
	}
	return
}

func (d Dealer) score() (total int) {
	for _, c := range d.hand {
		total += scoreCard(c)
	}
	total += scoreCard(d.hiddenCard)
	if d.hasAce() && total <= 11 {
		total += 10
	}
	return
}

func (p Player) hasAce() bool {
	for _, c := range p.hand {
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
	score = int(c.Value) + 1
	if score > 10 && score <= 13 {
		score = 10
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
