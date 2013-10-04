package checkerlution

import (
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
)

type CheckerlutionScape struct {
	thinker *Checkerlution
}

func (scape *CheckerlutionScape) SetThinker(thinker *Checkerlution) {
	scape.thinker = thinker
}

func (scape *CheckerlutionScape) Fitness(cortex *ng.Cortex) (fitness float64) {

	logg.LogTo("MAIN", "scape.Fitness() called, playing checkers game")

	// play a game of checkers
	scape.thinker.StartWithCortex(cortex, cbot.RED_TEAM)
	game := cbot.NewGame(cbot.RED_TEAM, scape.thinker)
	game.GameLoop()
	logg.LogTo("MAIN", "gameLoop finished, get fitness")

	// get result
	fitness = scape.thinker.latestFitnessScore
	logg.LogTo("MAIN", "checkerlution scape returning fitness: %v", fitness)

	cortex.Shutdown()

	return

}
