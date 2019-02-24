package main

import (
	"fmt"

	"github.com/blogscot/deck"
)

func main() {
	fmt.Println("Let's play BlackJack!")
	fmt.Println("\nThe dealer shuffles the deck thoroughly then starts dealing...")
	fmt.Println()

	deck := deck.New()
	Start(&deck)

	justPlayersLength := len(players) - 1
	ps := players[:justPlayersLength]
	for _, p := range ps {
		play(&p)
	}

	d := players[justPlayersLength]
	play(&d)

	winner := decideWinner(players)
	if winner == "You" {
		fmt.Printf("\n%s win!\n", winner)
	} else if winner == gameIsDrawn {
		fmt.Printf("\nThe game is a draw!")
	} else {
		fmt.Printf("\n%s wins!\n", winner)
	}
}
