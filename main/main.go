package main

import (
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
)

func main() {

	logg.LogKeys["MAIN"] = true

	ng.SeedRandom()

	thinker := new(checkerlution.Checkerlution)
	redTeam := cbot.RED_TEAM
	game := cbot.NewGame(redTeam, thinker)
	game.GameLoop()

}
