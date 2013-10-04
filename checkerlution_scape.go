package checkerlution

import (
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
)

type CheckerlutionScape struct {
}

func (scape CheckerlutionScape) Fitness(cortex *ng.Cortex) (fitness float64) {

	// play a game of checkers
	thinker := &Checkerlution{mode: TRAINING_MODE}
	thinker.StartWithCortex(cortex, cbot.RED_TEAM)
	game := cbot.NewGame(cbot.RED_TEAM, thinker)
	game.GameLoop()

	// get result
	fitness = thinker.latestFitnessScore
	logg.LogTo("MAIN", "checkerlution scape returning fitness: %v", fitness)

	return

}
