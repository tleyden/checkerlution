package main

import (
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	ng "github.com/tleyden/neurgo"
	nv "github.com/tleyden/neurvolve"
)

func init() {
	logg.LogKeys["MAIN"] = true
	logg.LogKeys["NODE_SEND"] = true
	logg.LogKeys["NODE_RECV"] = true
	ng.SeedRandom()
}

func main() {

	// create a checkerlution instance just to create a cortex (kludgy)
	thinker := &checkerlution.Checkerlution{}
	thinker.SetMode(checkerlution.TRAINING_MODE)
	logg.LogTo("MAIN", "main checkerlution: %v", thinker)
	thinker.CreateNeurgoCortex()
	cortex := thinker.Cortex()
	logg.LogTo("MAIN", "main cortex: %p", cortex)

	// setup the scape
	scape := &checkerlution.CheckerlutionScape{}
	scape.SetThinker(thinker)

	// create a stochastic hill climber
	shc := &nv.StochasticHillClimber{
		FitnessThreshold:           50,
		MaxIterationsBeforeRestart: 5,
		MaxAttempts:                1,
	}
	cortexTrained, succeeded := shc.TrainScape(cortex, scape)
	if succeeded {
		logg.LogTo("MAIN", "Training succeeded, dumping to file")
		cortexTrained.MarshalJSONToFile("/tmp/checkerlution.json")
	} else {
		logg.LogTo("MAIN", "Training Failed")
	}

}
