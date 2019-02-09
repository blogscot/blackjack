package main

import (
	"fmt"
)

func main() {

	fmt.Println("Let's play BlackJack!")
	fmt.Println("\nThe dealer shuffles the deck thoroughly then startings dealing...\n")
	startGame()

	gameStatus()

	p := players[0]
	if playerChoice() == hit {
		p.takeCard(dealCard())
	}
	gameStatus()

}
