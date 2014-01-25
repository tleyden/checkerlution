package main

import (
	_ "expvar"
	"flag"
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	cbot "github.com/tleyden/checkers-bot"
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

func train(checkersBotFlags cbot.CheckersBotFlags) {

	// run a webserver in order to view expvar output
	// at http://localhost:8080/debug/vars
	go http.ListenAndServe(":8080", nil)

	trainer := &checkerlution.CheckerlutionTrainer{}

	// checkerlution.RunTopologyMutatingTrainer()
	trainer.RunPopulationTrainer(checkersBotFlags)

}

func run(checkersBotFlags cbot.CheckersBotFlags) {

	LOAD_CORTEX_FROM_FILE := false

	thinker := &checkerlution.Checkerlution{}
	thinker.SetMode(checkerlution.RUNNING_MODE)

	if LOAD_CORTEX_FROM_FILE {
		filename := "checkerlution_trained.json"
		cortex, err := ng.NewCortexFromJSONFile(filename)
		if err != nil {
			logg.LogPanic("Error reading cortex from: %v.  Err: %v", filename, err)
		}
		thinker.StartWithCortex(cortex, checkersBotFlags.Team)

	} else {
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

	var mode int

	flag.IntVar(
		&mode,
		"mode",
		0,
		"Run Mode - Run: 0 Train: 1",
	)

	flag.Parse()

	checkersBotFlags := checkersBotRawFlags.GetCheckersBotFlags()

	logg.LogTo("CHECKERLUTION", "Flags: %v", checkersBotFlags)
	logg.LogTo("CHECKERLUTION", "Mode: %v", mode)

	if mode == 0 {
		logg.LogTo("CHECKERLUTION", "Run mode: Run")
		run(checkersBotFlags)
	} else {
		logg.LogTo("CHECKERLUTION", "Run mode: Train")
		train(checkersBotFlags)
	}

}
