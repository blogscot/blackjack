package blackjack

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/blogscot/deck"
)

// Participant ...
type Participant interface {
	score() int
	add(c deck.Card)
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
			arr = append(arr, dealer.hiddenCard.String())
		} else {
			arr = append(arr, "and a card face down")
		}
	}

	return strings.Join(arr, ", ")
}

func isBust(s Participant) bool {
	return s.score() > 21
}

func play(s *Participant) (err error) {
	switch p := (*s).(type) {
	case *Player:
		err = handlePlayer(p)
	case *Dealer:
		err = handleDealer(p, player1.score())
	}
	return err
}

func handleDealer(d *Dealer, playerScore int) error {
	showHands(true)

	score := d.score()
	hasSoft17 := score == 7 && d.hasAce()
	if hasSoft17 {
		fmt.Println("The dealer has a soft 17.")
	}

	for score < 16 && score < playerScore || hasSoft17 {
		fmt.Print("The dealer hits, ")
		newCard := dealer.deal()
		d.add(newCard)
		if isBust(d) {
			return errors.New("The dealer is bust")
		}
	}

	fmt.Println("The dealer stands.")
	d.total = d.score()
	showHands(true)
	return nil
}

func handlePlayer(p *Player) error {
	giveMe := true

	for giveMe {
		if playerChoice() == hit {
			newCard := dealer.deal()
			p.add(newCard)
			if isBust(p) {
				return fmt.Errorf("%s", p.getName())
			}
			showHands(false)
		} else {
			giveMe = false
		}
	}
	p.total = p.score()
	return nil
}

func decideWinner(s []Participant) string {
	sort.Slice(s, func(i, j int) bool {
		return s[i].score() > s[j].score()
	})
	if s[0].score() != s[1].score() {
		return s[0].getName()
	}
	return gameIsDrawn
}
