package main

import (
	"fmt"
)

func main() {

	fmt.Println("Let's play BlackJack!\n")
	fmt.Println("The dealer shuffles the deck thoroughly then starts dealing...\n")
	startGame()

	p := players[0]
	play(&p)

	d := players[len(players)-1]
	play(&d)

}
