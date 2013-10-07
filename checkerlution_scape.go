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

	cortex.Init()

	// TODO: if this is the 2nd or later round, and Fitness is called
	// even though the residue from the previous round is still hanging
	// around .. it will think the game is finished and return a fitness.
	// we basically need to wait/poll until the game number changes.
	// and we'll probably have to store the game number in the scape.
	// OR a simpler approach - the GameLoop won't unblock until
	// the next game is ready

	// play a game of checkers
	scape.thinker.StartWithCortex(cortex, cbot.RED_TEAM)
	game := cbot.NewGame(cbot.RED_TEAM, scape.thinker)
	game.GameLoop()
	logg.LogTo("MAIN", "gameLoop finished")

	// get result (TODO: this feels clunky, just call thinker.calcFitness() here)
	fitness = scape.thinker.latestFitnessScore
	logg.LogTo("MAIN", "checkerlution scape fitness: %v", fitness)

	// wait until the game number increments, otherwise on the
	// next callback to this method, we'll jump into a game which
	// is already over.
	logg.LogTo("MAIN", "checkerlution WaitForNextGame()")
	game.WaitForNextGame()
	logg.LogTo("MAIN", "/ finished checkerlution WaitForNextGame()")

	// re-initialize the thinker (TODO: add method for this)
	cortex.Shutdown()

	return

}
