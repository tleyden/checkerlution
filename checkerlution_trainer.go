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

func (trainer *CheckerlutionTrainer) RunPopulationTrainer(checkerlutionFlags CheckerlutionFlags) {

	checkersBotFlags := checkerlutionFlags.CheckersBotFlags

	// setup the scape
	scape := &CheckerlutionScape{}
	scape.SetSyncGatewayUrl(checkersBotFlags.SyncGatewayUrl)
	scape.SetFeedType(checkersBotFlags.FeedType)
	scape.SetTeam(checkersBotFlags.Team)
	scape.SetRandomDelayBeforeMove(checkersBotFlags.RandomDelayBeforeMove)

	// create population trainer ...
	pt := &nv.PopulationTrainer{
		FitnessThreshold: 1000, // set very high, will never hit it ..
		MaxGenerations:   25,
		CortexMutator:    nv.MutateAllWeightsBellCurve,
		NumOpponents:     5,
	}

	nv.RegisterHandlers(pt)

	recorder := NewRecorder(
		checkersBotFlags.SyncGatewayUrl,
		checkerlutionFlags.PopulationName,
	)

	cortexes := []*ng.Cortex{}
	if len(recorder.GetLatestGenerationCortexes()) > 0 {
		cortexes = recorder.GetLatestGenerationCortexes()
		if len(cortexes) == 0 {
			logg.LogPanic("No cortexes found in latest generation")
		}
		logg.LogTo("CHECKERLUTION", "Starting from existing generation")

	} else {
		cortexes = getInitialGeneration()
		logg.LogTo("CHECKERLUTION", "Starting from new")
	}

	for _, cortex := range cortexes {
		logg.LogTo("CHECKERLUTION", "Cortex: %v", cortex.NodeId.UUID)
	}

	fitGeneration, succeeded := pt.Train(cortexes, scape, recorder)

	if succeeded {
		logg.LogTo("MAIN", "Training finished (exceeded threshold)")
	} else {
		logg.LogTo("MAIN", "Training finished (did not exceed threshold)")
	}

	for i, evaldCortex := range fitGeneration {

		logg.LogTo("MAIN", "Cortex %d fitness: %v", i, evaldCortex.Fitness)

	}

}

func (trainer *CheckerlutionTrainer) RunTopologyMutatingTrainer(checkersBotFlags cbot.CheckersBotFlags) {

	// setup the scape
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
	for i := 0; i < 30; i++ {

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
