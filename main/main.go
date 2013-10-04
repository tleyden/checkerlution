package main

import (
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	cbot "github.com/tleyden/checkers-bot"
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
	cortex := thinker.cortex

	// setup the scape
	scape := NewCheckerlutionScape()

	// create a stochastic hill climber
	shc := &nv.StochasticHillClimber{
		FitnessThreshold:           ng.FITNESS_THRESHOLD,
		MaxIterationsBeforeRestart: 100,
		MaxAttempts:                10,
	}
	cortexTrained, succeeded := shc.TrainScape(cortex)

}
