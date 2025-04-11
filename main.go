package main

import (
	"log"
)

func main() {
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	defer game.screen.Fini()

	game.Run()
}
