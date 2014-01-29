package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	cbot "github.com/tleyden/checkers-bot"
	"github.com/tleyden/go-couch"
	ng "github.com/tleyden/neurgo"
	"net/http"
)

func init() {
	logg.LogKeys["MAIN"] = true
	logg.LogKeys["DEBUG"] = true
	logg.LogKeys["CHECKERLUTION"] = true
	logg.LogKeys["CHECKERLUTION_SCAPE"] = true
	logg.LogKeys["CHECKERSBOT"] = true
	logg.LogKeys["NEURGO"] = false
	logg.LogKeys["NEURVOLVE"] = true
	logg.LogKeys["SENSOR_SYNC"] = false
	logg.LogKeys["ACTUATOR_SYNC"] = false
	logg.LogKeys["NODE_PRE_SEND"] = false
	logg.LogKeys["NODE_POST_SEND"] = false
	logg.LogKeys["NODE_POST_RECV"] = false
	logg.LogKeys["NODE_STATE"] = false
	ng.SeedRandom()
}

func train(checkerlutionFlags checkerlution.CheckerlutionFlags) {

	// run a webserver in order to view expvar output
	// at http://localhost:8080/debug/vars
	go http.ListenAndServe(":8080", nil)

	trainer := &checkerlution.CheckerlutionTrainer{}

	// checkerlution.RunTopologyMutatingTrainer()
	trainer.RunPopulationTrainer(checkerlutionFlags)

}

func run(checkersBotFlags cbot.CheckersBotFlags, cortexId string) {

	thinker := &checkerlution.Checkerlution{}
	thinker.SetMode(checkerlution.RUNNING_MODE)

	if len(cortexId) > 0 {
		// first try to load it from a file with that id
		filename := fmt.Sprintf("%v.json", cortexId)
		cortex, err := ng.NewCortexFromJSONFile(filename)
		if err == nil {
			thinker.StartWithCortex(cortex, checkersBotFlags.Team)

		} else {
			// otherwise, load it from db
			db, error := couch.Connect(checkersBotFlags.SyncGatewayUrl)
			if error != nil {
				logg.LogPanic("Error connecting to %v: %v", checkersBotFlags.SyncGatewayUrl, error)
			}

			cortex := &ng.Cortex{}

			error = db.Retrieve(cortexId, cortex)
			logg.LogTo("CHECKERLUTION", "error: %v cortex :%v", error, *cortex)

			if error != nil {
				logg.LogPanic("Could not find cortex: %v", cortexId, error)
			}

			cortex.LinkNodesToCortex()
			thinker.StartWithCortex(cortex, checkersBotFlags.Team)

		}

	} else {
		// start with a random cortex
		thinker.Start(checkersBotFlags.Team)
	}

	game := cbot.NewGame(checkersBotFlags.Team, thinker)
	game.SetServerUrl(checkersBotFlags.SyncGatewayUrl)
	game.SetFeedType(checkersBotFlags.FeedType)
	game.SetDelayBeforeMove(checkersBotFlags.RandomDelayBeforeMove)

	logg.LogTo("CHECKERLUTION", "Starting game loop")
	game.GameLoop()
	logg.LogTo("CHECKERLUTION", "Game loop finished")

}

func main() {

	checkersBotRawFlags := cbot.GetCheckersBotRawFlags()

	// TODO: make this a friendly string instead of an int
	var mode int

	flag.IntVar(
		&mode,
		"mode",
		0,
		"Run Mode - Run: 0 Train: 1",
	)

	var cortexId string

	flag.StringVar(
		&cortexId,
		"cortexId",
		"",
		"Cortex Id of saved cortex",
	)

	checkerlutionFlags := checkerlution.CheckerlutionFlags{}
	flag.StringVar(
		&checkerlutionFlags.PopulationName,
		"populationName",
		"",
		"Population name, eg, population1",
	)

	flag.Parse()

	checkersBotFlags := checkersBotRawFlags.GetCheckersBotFlags()
	checkerlutionFlags.CheckersBotFlags = checkersBotFlags

	logg.LogTo("CHECKERLUTION", "Flags: %v", checkersBotFlags)
	logg.LogTo("CHECKERLUTION", "Mode: %v", mode)
	logg.LogTo("CHECKERLUTION", "CortexId: %v", cortexId)

	if mode == 0 {
		logg.LogTo("CHECKERLUTION", "Run mode: Run")
		run(checkersBotFlags, cortexId)
	} else {
		logg.LogTo("CHECKERLUTION", "Run mode: Train")
		if len(checkerlutionFlags.PopulationName) == 0 {
			logg.LogPanic("populationName required in training mode")
		}

		train(checkerlutionFlags)
	}

}
