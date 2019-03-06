package blackjack

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	deck "github.com/blogscot/card-deck"
)

type Participant interface {
	score() int
	add(c deck.Card)
	getName() string
}

// Show the player and dealers hands. The dealer keeps one card
// hidden until the player has finished their turn.
func showHands(hideCard bool) {
	fmt.Print(pageBreak)
	fmt.Print("Player1 has: ")
	fmt.Printf("%s. (score %d)", getPlayerHand(player1), player1.score())

	fmt.Print("\nThe dealer has: ")
	dealerCards := fmt.Sprintf("%s.", getDealerHand(dealer, hideCard))
	if !hideCard {
		dealerCards += fmt.Sprintf(" (score %d)\n", dealer.score())
	}
	fmt.Print(dealerCards)
}

func getPlayerHand(p Player) string {
	arr := []string{}
	for _, card := range player1.hand {
		arr = append(arr, card.String())
	}
	return strings.Join(arr, ", ")
}

func getDealerHand(d Dealer, hideCard bool) string {
	arr := []string{}
	cards := d.hand

	if hideCard {
		cards = cards[:len(cards)-1]
	}
	for _, card := range cards {
		arr = append(arr, card.String())
	}
	if hideCard {
		arr = append(arr, "and a hidden card.\n")
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
	return
}

func handleDealer(d *Dealer, playerScore int) error {
	showHands(false)

	score := d.score()
	if score > playerScore {
		return dealerStands()
	}

	hasSoft17 := score == 17 && d.hasAce()
	for score < playerScore || hasSoft17 {
		fmt.Print("The dealer hits, ")
		newCard := d.dealCard()
		d.add(newCard)
		fmt.Printf("and receives %s. (score %d)\n", newCard, d.score())
		score = d.score()
		if isBust(d) {
			return errors.New("The dealer is bust")
		}
		if score > playerScore || score == 21 {
			return dealerStands()
		}
	}

	return dealerStands()
}

func dealerStands() error {
	fmt.Println("The dealer stands.")
	showHands(false)
	return nil
}

func handlePlayer(p *Player) error {
	giveMe := true

	for giveMe {
		if playerChoice() == hit {
			newCard := dealer.dealCard()
			fmt.Printf("You receive %s.\n", newCard)
			p.add(newCard)
			if isBust(p) {
				return fmt.Errorf("%s", p.getName())
			}
			showHands(true)
		} else {
			giveMe = false
		}
	}
	return nil
}

// decideWinner find the winning player.
// It assumes that players haven't gone bust already.
func decideWinner(s []Participant) string {
	sort.Slice(s, func(i, j int) bool {
		return s[i].score() > s[j].score()
	})
	if s[0].score() != s[1].score() {
		return s[0].getName()
	}
	return gameIsDrawn
}
