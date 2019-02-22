package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/blogscot/deck"
)

// Scoring ...
type Scoring interface {
	score() int
	takeCard(c deck.Card)
	getName() string
}

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
)

var (
	reader = bufio.NewReader(os.Stdin)

	dealer  = Dealer{Player: Player{name: "The dealer"}}
	player1 = Player{name: "You"}
	players = []Scoring{&player1, &dealer}
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
	return
}

func (d Dealer) score() (total int) {
	return d.Player.score() + scoreCard(d.faceDown)
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

func showHands(showAll bool) {
	fmt.Print(pageBreak)
	fmt.Print("You have: ")
	fmt.Printf("%s.", showHand(&player1, false))

	fmt.Print("\nThe dealer has: ")
	fmt.Printf("%s.\n", showHand(&dealer, showAll))
}

func showHand(s Scoring, showAll bool) string {
	arr := []string{}

	switch s.(type) {
	case *Player:
		for _, card := range player1.cards {
			arr = append(arr, card.String())
		}
	case *Dealer:
		for _, card := range dealer.cards {
			arr = append(arr, card.String())
		}
		if showAll {
			arr = append(arr, dealer.faceDown.String())
		} else {
			arr = append(arr, "and a card face down")
		}
	}

	return strings.Join(arr, ", ")
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

func isBust(s Scoring) bool {
	return s.score() > 21
}

func play(s *Scoring) {
	switch p := (*s).(type) {
	case *Player:
		handlePlayer(p)
	case *Dealer:
		handleDealer(p)
	}
}

func hasAce(cs []deck.Card) bool {
	for _, c := range cs {
		if c.Value == deck.Ace {
			return true
		}
	}
	return false
}

func handleDealer(d *Dealer) {
	showHands(true)

	score := d.score()
	hasSoft17 := score == 7 && hasAce(d.cards)
	if hasSoft17 {
		fmt.Println("The dealer has a soft 17.")
	}

	for score < 16 || hasSoft17 {
		d.takeCard(dealCard())
		if isBust(d) {
			fmt.Println(dealerIsBust)
			os.Exit(0)
		}
	}
	d.total = d.score()
	showHands(true)
}

func handlePlayer(p *Player) {
	giveMe := true

	for giveMe {
		if playerChoice() == hit {
			p.takeCard(dealCard())
			if isBust(p) {
				fmt.Println(playerIsBust)
				os.Exit(0)
			}
			showHands(false)
		} else {
			giveMe = false
		}
	}
	p.total = p.score()
}

func decideWinner(s []Scoring) string {
	sort.Slice(s, func(i, j int) bool {
		return s[i].score() > s[j].score()
	})
	if s[0].score() != s[1].score() {
		return s[0].getName()
	}
	return dealer.name
}
