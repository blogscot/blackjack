package main

import (
	"github.com/blogscot/blackjack"
	"github.com/blogscot/deck"
)

func main() {
	deck := deck.New()

	blackjack.Play(&deck)
}
