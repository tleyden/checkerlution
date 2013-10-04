package main

import (
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	ng "github.com/tleyden/neurgo"
	nv "github.com/tleyden/neurvolve"
)

func init() {
	logg.LogKeys["MAIN"] = true
	ng.SeedRandom()
}

func main() {

	// create a checkerlution instance just to create a cortex (kludgy)
	thinker := new(checkerlution.Checkerlution)
	thinker.CreateNeurgoCortex()
	cortex := thinker.Cortex()

	// setup the scape
	scape := &checkerlution.CheckerlutionScape{}

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
