package main

import (
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	ng "github.com/tleyden/neurgo"
)

func main() {

	logg.LogKeys["MAIN"] = true

	ng.SeedRandom()

	redTeam := checkerlution.RED_TEAM
	game := checkerlution.NewGame(redTeam)
	game.GameLoop()

}
