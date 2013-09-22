package main

import (
	"github.com/tleyden/checkerlution"
	ng "github.com/tleyden/neurgo"
	"log"
)

func main() {

	log.Printf("hello")

	ng.SeedRandom()

	redTeam := checkerlution.RED_TEAM
	game := checkerlution.NewGame(redTeam)
	game.GameLoop()

}
