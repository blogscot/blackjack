package main

import (
	"fmt"
)

func main() {
	fmt.Println("Let's play BlackJack!\n")
	fmt.Println("The dealer shuffles the deck thoroughly then starts dealing...\n")
	startGame()

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
	} else {
		fmt.Printf("\n%s wins!\n", winner)
	}
}
