package main

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
	nv "github.com/tleyden/neurvolve"
	"math"
	"time"
)

func init() {
	logg.LogKeys["MAIN"] = true
	logg.LogKeys["DEBUG"] = false
	logg.LogKeys["NEURGO"] = false
	logg.LogKeys["NODE_PRE_SEND"] = false
	logg.LogKeys["NODE_POST_SEND"] = false
	logg.LogKeys["NODE_POST_RECV"] = false
	logg.LogKeys["NODE_STATE"] = false
	ng.SeedRandom()
}

func main() {

	// create a checkerlution instance just to create a cortex (kludgy)
	thinker := &checkerlution.Checkerlution{}
	thinker.SetMode(checkerlution.TRAINING_MODE)
	thinker.CreateNeurgoCortex()
	// thinker.LoadNeurgoCortex("/tmp/checkerlution-good.json")
	cortex := thinker.Cortex()

	// setup the scape
	checkersBotFlags := cbot.ParseCmdLine()
	scape := &checkerlution.CheckerlutionScape{}
	scape.SetThinker(thinker)
	scape.SetSyncGatewayUrl(checkersBotFlags.SyncGatewayUrl)
	scape.SetFeedType(checkersBotFlags.FeedType)
	scape.SetTeam(checkersBotFlags.Team)
	scape.SetRandomDelayBeforeMove(checkersBotFlags.RandomDelayBeforeMove)

	// create a stochastic hill climber
	shc := &nv.StochasticHillClimber{
		FitnessThreshold:           150,
		MaxIterationsBeforeRestart: 10,
		MaxAttempts:                500,
		WeightSaturationRange:      []float64{-2 * math.Pi, 2 * math.Pi},
	}
	tmt := &nv.TopologyMutatingTrainer{
		MaxAttempts:                100,
		MaxIterationsBeforeRestart: 5,
		StochasticHillClimber:      shc,
	}
	cortexTrained, succeeded := tmt.Train(cortex, scape)
	if succeeded {
		logg.LogTo("MAIN", "Training succeeded")
	} else {
		logg.LogTo("MAIN", "Training Failed")
	}

	filename := fmt.Sprintf("/tmp/checkerlution-%v.json", time.Now().Unix())
	logg.LogTo("MAIN", "Saving Cortex to %v", filename)
	cortexTrained.MarshalJSONToFile(filename)

}
