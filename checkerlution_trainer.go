package checkerlution

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
	nv "github.com/tleyden/neurvolve"
	"math"
	"time"
)

func runPopulationTrainer() {

	// setup the scape
	checkersBotFlags := cbot.ParseCmdLine()
	scape := &CheckerlutionScapeTwoPlayer{}
	scape.SetSyncGatewayUrl(checkersBotFlags.SyncGatewayUrl)
	scape.SetFeedType(checkersBotFlags.FeedType)
	scape.SetTeam(checkersBotFlags.Team)
	scape.SetRandomDelayBeforeMove(checkersBotFlags.RandomDelayBeforeMove)

	noOpCortexMutator := func(cortex *ng.Cortex) (success bool, result nv.MutateResult) {
		success = true
		result = "nothing"
		return
	}

	// create population trainer ...
	pt := &nv.PopulationTrainer{
		FitnessThreshold: 150,
		MaxGenerations:   1,
		CortexMutator:    noOpCortexMutator,
	}

	population := getInitialPopulation()

	fitPopulation, succeeded := pt.Train(population, scape)

	if succeeded {
		logg.LogTo("MAIN", "Training succeeded")
	} else {
		logg.LogTo("MAIN", "Training Failed")
	}

	for i, fitCortex := range fitPopulation {

		logg.LogTo("MAIN", "Cortex %d fitness: %v", i, fitCortex.Fitness)
		filename := fmt.Sprintf("/tmp/checkerlution-%v.json", time.Now().Unix())
		logg.LogTo("MAIN", "Saving Cortex to %v", filename)
		cortex := fitCortex.Cortex
		cortex.MarshalJSONToFile(filename)

	}

}

func runTopologyMutatingTrainer() {

	// create a checkerlution instance just to create a cortex (kludgy)
	thinker := &Checkerlution{}
	thinker.SetMode(TRAINING_MODE)
	// thinker.CreateNeurgoCortex()
	thinker.LoadNeurgoCortex("cortex_avg8.json")
	cortex := thinker.Cortex()

	// setup the scape
	checkersBotFlags := cbot.ParseCmdLine()
	scape := &CheckerlutionScape{}
	scape.SetThinker(thinker)
	scape.SetSyncGatewayUrl(checkersBotFlags.SyncGatewayUrl)
	scape.SetFeedType(checkersBotFlags.FeedType)
	scape.SetTeam(checkersBotFlags.Team)
	scape.SetRandomDelayBeforeMove(checkersBotFlags.RandomDelayBeforeMove)

	// create a stochastic hill climber (required by topology mutation trainer)
	shc := &nv.StochasticHillClimber{
		FitnessThreshold:           150,
		MaxIterationsBeforeRestart: 10,
		MaxAttempts:                500,
		WeightSaturationRange:      []float64{-2 * math.Pi, 2 * math.Pi},
	}

	// this thing will train the network by randomly mutating and calculating fitness
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

func getInitialPopulation() (population []*ng.Cortex) {
	population = make([]*ng.Cortex, 0)
	for i := 0; i < 30; i++ {

		thinker := &Checkerlution{}

		// probably a bug, because the cortex is now "bound" to this
		// thinker .. but we're discarding the thinker
		thinker.CreateNeurgoCortex()

		population = append(population, thinker.Cortex())
	}
	return
}
