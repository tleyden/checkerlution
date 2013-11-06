package checkerlution

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
	"time"
)

type CheckerlutionScape struct {
	team                  cbot.TeamType
	syncGatewayUrl        string
	feedType              cbot.FeedType
	randomDelayBeforeMove int
}

func (scape *CheckerlutionScape) Fitness(cortex *ng.Cortex) (fitness float64) {

	logg.LogTo("DEBUG", "scape.Fitness() called, playing checkers game")

	fitnessVals := []float64{}

	for i := 0; i < 7; i++ {

		cortex.Init()

		// create a checkerlution instance
		thinker := &Checkerlution{}
		thinker.SetMode(TRAINING_MODE)
		thinker.StartWithCortex(cortex, scape.team)

		game := cbot.NewGame(scape.team, thinker)
		game.SetServerUrl(scape.syncGatewayUrl)
		game.SetFeedType(scape.feedType)
		game.SetDelayBeforeMove(scape.randomDelayBeforeMove)
		game.GameLoop()
		logg.LogTo("DEBUG", "gameLoop finished")

		latestFitness := thinker.latestFitnessScore
		logg.LogTo("MAIN", "Fitness: %v", latestFitness)
		fitnessVals = append(fitnessVals, latestFitness)

		// wait until the game number increments, otherwise on the
		// next callback to this method, we'll jump into a game which
		// is already over.
		logg.LogTo("MAIN", "Wait For Next Game ..")
		game.WaitForNextGame()

		cortex.Shutdown()

	}

	logg.LogTo("DEBUG", "Game Set Finished")

	total := 0.0
	for _, fitnessVal := range fitnessVals {
		total += fitnessVal
	}
	fitness = total / float64(len(fitnessVals))
	logg.LogTo("MAIN", "Average Fitness for Game Set: %v", fitness)

	filename := fmt.Sprintf("/tmp/checkerlution-%v.json", time.Now().Unix())
	logg.LogTo("MAIN", "Saving Cortex snapshot to %v", filename)
	cortex.MarshalJSONToFile(filename)

	return

}

func (scape *CheckerlutionScape) FitnessAgainst(cortex *ng.Cortex, opponentCortex *ng.Cortex) (fitness float64) {
	logg.LogPanic("Not implemented")
	fitness = 0.0
	return
}

func (scape *CheckerlutionScape) SetSyncGatewayUrl(syncGatewayUrl string) {
	scape.syncGatewayUrl = syncGatewayUrl
}

func (scape *CheckerlutionScape) SetTeam(team cbot.TeamType) {
	scape.team = team
}

func (scape *CheckerlutionScape) SetFeedType(feedType cbot.FeedType) {
	scape.feedType = feedType
}

func (scape *CheckerlutionScape) SetRandomDelayBeforeMove(delay int) {
	scape.randomDelayBeforeMove = delay
}
