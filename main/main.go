package main

import (
	"github.com/blogscot/blackjack"
)

func main() {
	deck := blackjack.Deck()
	blackjack.Play(&deck)
}
