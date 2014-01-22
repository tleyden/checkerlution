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

type CheckerlutionTrainer struct{}

func (trainer *CheckerlutionTrainer) RunPopulationTrainer() {

	// setup the scape
	checkersBotFlags := cbot.ParseCmdLine()
	scape := &CheckerlutionScape{}
	scape.SetSyncGatewayUrl(checkersBotFlags.SyncGatewayUrl)
	scape.SetFeedType(checkersBotFlags.FeedType)
	scape.SetTeam(checkersBotFlags.Team)
	scape.SetRandomDelayBeforeMove(checkersBotFlags.RandomDelayBeforeMove)

	// create population trainer ...
	pt := &nv.PopulationTrainer{
		FitnessThreshold: 1000, // set very high, will never hit it ..
		MaxGenerations:   10,
		CortexMutator:    nv.MutateAllWeightsBellCurve,
		NumOpponents:     3,
	}

	generation := getInitialGeneration()
	nv.RegisterHandlers(pt)

	population := Population{name: "population11"}
	recorder := NewRecorder(checkersBotFlags.SyncGatewayUrl, population)

	fitGeneration, succeeded := pt.Train(generation, scape, recorder)

	if succeeded {
		logg.LogTo("MAIN", "Training finished (exceeded threshold)")
	} else {
		logg.LogTo("MAIN", "Training finished (did not exceed threshold)")
	}

	for i, evaldCortex := range fitGeneration {

		logg.LogTo("MAIN", "Cortex %d fitness: %v", i, evaldCortex.Fitness)
		filename := fmt.Sprintf("/tmp/checkerlution-%v.json", time.Now().Unix())
		logg.LogTo("MAIN", "Saving Cortex to %v", filename)
		cortex := evaldCortex.Cortex
		cortex.MarshalJSONToFile(filename)

	}

}

func (trainer *CheckerlutionTrainer) RunTopologyMutatingTrainer() {

	// setup the scape
	checkersBotFlags := cbot.ParseCmdLine()
	scape := &CheckerlutionScape{}
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
	initialCortex := getInitialCortex()
	cortexTrained, succeeded := tmt.Train(initialCortex, scape)
	if succeeded {
		logg.LogTo("MAIN", "Training succeeded")
	} else {
		logg.LogTo("MAIN", "Training Failed")
	}

	filename := fmt.Sprintf("/tmp/checkerlution-%v.json", time.Now().Unix())
	logg.LogTo("MAIN", "Saving Cortex to %v", filename)
	cortexTrained.MarshalJSONToFile(filename)

}

func getInitialGeneration() (generation []*ng.Cortex) {
	generation = make([]*ng.Cortex, 0)
	for i := 0; i < 10; i++ {

		thinker := &Checkerlution{}

		thinker.CreateNeurgoCortex()

		generation = append(generation, thinker.Cortex())
	}
	return
}

func getInitialCortex() (cortex *ng.Cortex) {
	thinker := &Checkerlution{}
	thinker.CreateNeurgoCortex()
	return thinker.Cortex()
}
