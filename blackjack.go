package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blogscot/deck"
)

// Player ...
type Player struct {
	name  string
	cards []deck.Card
	total int
}

// Dealer behaves like a player but keeps one card hidden
type Dealer struct {
	Player
	faceDown deck.Card
}

type choice int

const (
	hit = iota
	stand
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
)

var (
	reader = bufio.NewReader(os.Stdin)

	dealer  = Dealer{Player: Player{name: "The dealer"}}
	player1 = Player{name: "You"}
	players = []Participant{&player1, &dealer}
	cards   = deck.New()
)

func startGame() {
	cards.Shuffle()

	for n := 1; n <= 2; n++ {
		for _, p := range players {
			c := dealCard()
			p.takeCard(c)
		}
	}
	showHands(false)
}

func dealCard() (card deck.Card) {
	card = cards[0]
	cards = cards[1:]
	return
}

func (p *Player) getName() string {
	return p.name
}

func (p *Player) takeCard(c deck.Card) {
	fmt.Printf("%s receive %s\n", p.name, c)
	p.cards = append(p.cards, c)
}

func (d *Dealer) takeCard(c deck.Card) {
	if d.faceDown == (deck.Card{}) {
		fmt.Printf("The dealer places a card face down.\n")
		d.faceDown = c
	} else {
		fmt.Printf("%s receives %s\n", d.name, c)
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
	total = d.Player.score() + scoreCard(d.faceDown)
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
	return d.Player.hasAce() || d.faceDown.Value == deck.Ace
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
