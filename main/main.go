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
	logg.LogKeys["DEBUG"] = true
	logg.LogKeys["NODE_PRE_SEND"] = true
	logg.LogKeys["NODE_POST_SEND"] = true
	logg.LogKeys["NODE_POST_RECV"] = true
	logg.LogKeys["NODE_STATE"] = true
	ng.SeedRandom()
}

func main() {

	// create a checkerlution instance just to create a cortex (kludgy)
	thinker := &checkerlution.Checkerlution{}
	thinker.SetMode(checkerlution.TRAINING_MODE)
	thinker.CreateNeurgoCortex()
	// thinker.LoadNeurgoCortex("/Users/traun/tmp/checkerlution-1381908895.json")
	cortex := thinker.Cortex()

	// setup the scape
	team, syncGatewayUrl, feedType := cbot.ParseCmdLine()
	scape := &checkerlution.CheckerlutionScape{}
	scape.SetThinker(thinker)
	scape.SetSyncGatewayUrl(syncGatewayUrl)
	scape.SetFeedType(feedType)
	scape.SetTeam(team)

	// create a stochastic hill climber
	shc := &nv.StochasticHillClimber{
		FitnessThreshold:           150,
		MaxIterationsBeforeRestart: 5,
		MaxAttempts:                5,
		WeightSaturationRange:      []float64{-2 * math.Pi, 2 * math.Pi},
	}
	cortexTrained, succeeded := shc.Train(cortex, scape)
	if succeeded {
		logg.LogTo("MAIN", "Training succeeded")
	} else {
		logg.LogTo("MAIN", "Training Failed")
	}

	filename := fmt.Sprintf("/tmp/checkerlution-%v.json", time.Now().Unix())
	logg.LogTo("MAIN", "Dumping latest cortex to %v", filename)
	cortexTrained.MarshalJSONToFile(filename)

}
