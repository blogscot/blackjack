package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/blogscot/deck"
)

// Participant ...
type Participant interface {
	score() int
	takeCard(c deck.Card)
	getName() string
}

func showHands(showAll bool) {
	fmt.Print(pageBreak)
	fmt.Print("You have: ")
	fmt.Printf("%s (score %d).", showHand(&player1, false), player1.score())

	fmt.Print("\nThe dealer has: ")
	fmt.Printf("%s.\n", showHand(&dealer, showAll))
}

func showHand(s Participant, showAll bool) string {
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

func isBust(s Participant) bool {
	return s.score() > 21
}

func play(s *Participant) {
	switch p := (*s).(type) {
	case *Player:
		handlePlayer(p)
	case *Dealer:
		handleDealer(p)
	}
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

func decideWinner(s []Participant) string {
	sort.Slice(s, func(i, j int) bool {
		return s[i].score() > s[j].score()
	})
	if s[0].score() != s[1].score() {
		return s[0].getName()
	}
	return dealer.name
}