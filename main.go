package main

import (
	"flag"
	"log"
)

func main() {
	langFlag := flag.String("lang", "", "Programming language (go/javascript/rust)")
	flag.Parse()

	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	defer game.screen.Fini()

	if *langFlag != "" {
		game.sentenceGen.SetLanguage(*langFlag)
		game.state = StatePlaying
	}

	game.Run()
}
